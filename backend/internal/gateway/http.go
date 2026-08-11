package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/contextx"
	"backend/internal/domain/auth"
	flowdomain "backend/internal/domain/flow"
	"backend/internal/permission"
	"backend/internal/repo"
	"backend/internal/secrets"
)

type HTTPGateway struct {
	config     config.GatewayConf
	auth       *auth.Service
	permission *permission.Evaluator
	routes     *repo.GatewayRepo
	flow       *flowdomain.Service
}

type runtimeConfig struct{}

const internalProxyTokenHeader = "X-Tier0-Internal-Token"

func NewHTTPGateway(ctx context.Context, c config.GatewayConf, authSvc *auth.Service, perm *permission.Evaluator, flowSvc *flowdomain.Service) *HTTPGateway {
	return &HTTPGateway{
		config:     c,
		auth:       authSvc,
		permission: perm,
		routes:     repo.NewGatewayRepo(ctx),
		flow:       flowSvc,
	}
}

func (g *HTTPGateway) Handle(w http.ResponseWriter, r *http.Request) bool {
	return g.handle(w, r)
}

func (g *HTTPGateway) handle(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == "/runtime-env.js" {
		g.serveRuntimeEnv(w, r)
		return true
	}
	if r.Method == http.MethodGet && (r.URL.Path == "/api/health" || r.URL.Path == "/api/health/") {
		g.serveHealth(w)
		return true
	}
	if isInternalPath(r.URL.Path) {
		return false
	}
	if strings.Contains(r.URL.Path, "..") {
		http.Error(w, "path traversal denied", http.StatusBadRequest)
		return true
	}
	// S3 object-storage prefixes are proxied before the configured gateway
	// routes: presigned URLs rewritten to the platform entrance must reach the
	// object store with the signing host/path intact for SigV4 verification.
	if g.tryGatewayRoute(w, r) {
		return true
	}
	return false
}

func (g *HTTPGateway) serveRuntimeEnv(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	payload, err := json.Marshal(runtimeConfig{})
	if err != nil {
		http.Error(w, "runtime config unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if r.Method == http.MethodHead {
		return
	}
	_, _ = w.Write([]byte("window.__TIER0_RUNTIME_CONFIG__ = "))
	_, _ = w.Write(payload)
	_, _ = w.Write([]byte(";\n"))
}

func (g *HTTPGateway) serveNodeRedFlowPayload(w http.ResponseWriter, r *http.Request, route repo.GatewayRoute) bool {
	flowType, idParam, cookieName, ok := nodeRedFlowPayloadRoute(r.URL.Path)
	if !ok {
		http.Error(w, "unsupported flow payload route", http.StatusBadGateway)
		return true
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return true
	}
	if !g.authorizeRoute(w, r, route) {
		return true
	}

	flowID := strings.TrimSpace(r.URL.Query().Get(idParam))
	if flowID == "" {
		flowID = flowIDFromReferer(r.Referer(), idParam)
	}
	if flowID == "" {
		if c, err := r.Cookie(cookieName); err == nil {
			flowID = strings.TrimSpace(c.Value)
		}
	}
	if flowID == "" {
		writeNodeRedPayload(w, flowdomain.EmptyEditorPayload())
		return true
	}

	id, err := strconv.ParseInt(flowID, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid flow id", http.StatusBadRequest)
		return true
	}
	if g.flow == nil {
		http.Error(w, "flow service unavailable", http.StatusServiceUnavailable)
		return true
	}
	payload, err := g.flow.EditorPayload(r.Context(), int64(flowType), id)
	if err != nil {
		if errors.Is(err, flowdomain.ErrNotFound) {
			http.Error(w, "flow not found", http.StatusNotFound)
			return true
		}
		if errors.Is(err, flowdomain.ErrInvalid) {
			http.Error(w, "invalid flow", http.StatusBadRequest)
			return true
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	writeNodeRedPayload(w, payload)
	return true
}

func nodeRedFlowPayloadRoute(path string) (int, string, string, bool) {
	parts := splitPathSegments(path)
	if len(parts) != 2 && len(parts) != 3 {
		return 0, "", "", false
	}
	if len(parts) == 2 && parts[1] != "flows" {
		return 0, "", "", false
	}
	if len(parts) == 3 && (parts[1] != "home" || parts[2] != "flows") {
		return 0, "", "", false
	}
	switch parts[0] {
	case "nodered":
		return 1, "flowId", "flowId", true
	case "eventflow":
		return 2, "flowId", "flowId", true
	default:
		return 0, "", "", false
	}
}

func flowIDFromReferer(referer, param string) string {
	if strings.TrimSpace(referer) == "" {
		return ""
	}
	u, err := url.Parse(referer)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(u.Query().Get(param))
}

func writeNodeRedPayload(w http.ResponseWriter, payload flowdomain.EditorPayload) {
	if payload.Flows == nil {
		payload.Flows = []map[string]any{}
	}
	if payload.Credentials == nil {
		payload.Credentials = map[string]map[string]string{}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{"flows": payload.Flows, "rev": payload.Rev, "credentials": payload.Credentials})
}

func (g *HTTPGateway) tryGatewayRoute(w http.ResponseWriter, r *http.Request) bool {
	routes, err := g.routes.ListGatewayRoutes(r.Context())
	if err != nil {
		http.Error(w, "gateway snapshot unavailable", http.StatusServiceUnavailable)
		return true
	}
	for _, route := range routes {
		if !route.Enabled {
			continue
		}
		if !routeMatches(route, r.Method, r.URL.Path) {
			continue
		}
		switch route.TargetType {
		case "reverseProxy":
			if !g.authorizeRoute(w, r, route) {
				return true
			}
			g.proxy(w, r, route)
			return true
		case "flowPayload":
			return g.serveNodeRedFlowPayload(w, r, route)
		case "static":
			if !g.authorizeRoute(w, r, route) {
				return true
			}
			g.serveStatic(w, r)
			return true
		default:
			http.Error(w, "unsupported gateway target type", http.StatusBadGateway)
			return true
		}
	}
	return false
}

func (g *HTTPGateway) authorizeRoute(w http.ResponseWriter, r *http.Request, route repo.GatewayRoute) bool {
	switch route.AuthPolicy {
	case "public":
		return true
	case "login":
		subject, err := g.auth.AuthenticateRequest(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return false
		}
		allowed, required, err := g.permission.Allow(r.Context(), subject, "gateway", r.Method, r.URL.Path)
		if err != nil {
			http.Error(w, "forbidden", http.StatusForbidden)
			return false
		}
		if allowed {
			return true
		}
		if required != "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return false
		}
		if route.ResourceKey != "" {
			if !contextx.HasResource(subject, route.ResourceKey) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return false
			}
			return true
		}
		http.Error(w, "forbidden", http.StatusForbidden)
		return false
	default:
		http.Error(w, "unsupported auth policy", http.StatusForbidden)
		return false
	}
}

func (g *HTTPGateway) proxy(w http.ResponseWriter, r *http.Request, route repo.GatewayRoute) {
	target, err := url.Parse(route.TargetUrl)
	if err != nil || target.Scheme == "" || target.Host == "" {
		http.Error(w, "invalid upstream", http.StatusBadGateway)
		return
	}
	if isDeniedTarget(target) {
		http.Error(w, "upstream denied by ssrf policy", http.StatusForbidden)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = rewriteProxyPath(r.URL.Path, route)
		req.URL.RawPath = ""
		req.Host = target.Host
		req.Header.Set("X-Forwarded-Host", r.Host)
		req.Header.Set("X-Forwarded-Proto", forwardedProto(r))
		injectInternalProxyAuth(req, route)
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		rewriteProxyLocation(resp, route, target)
		relaxEmbeddedProxyHeaders(resp, route)
		return nil
	}
	proxy.ServeHTTP(w, r)
}

func rewriteProxyLocation(resp *http.Response, route repo.GatewayRoute, target *url.URL) {
	location := strings.TrimSpace(resp.Header.Get("Location"))
	if location == "" {
		return
	}
	u, err := url.Parse(location)
	if err != nil {
		return
	}
	if u.IsAbs() {
		if !strings.EqualFold(u.Host, target.Host) {
			return
		}
		u.Scheme = ""
		u.Host = ""
	}
	prefix := strings.TrimSuffix(strings.TrimSuffix(route.PathPattern, "/**"), "/")
	if prefix == "" || !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, prefix+"/") || u.Path == prefix {
		resp.Header.Set("Location", u.String())
		return
	}
	u.Path = prefix + u.Path
	resp.Header.Set("Location", u.String())
}

func relaxEmbeddedProxyHeaders(resp *http.Response, route repo.GatewayRoute) {
	if strings.HasPrefix(route.PathPattern, "/emqx/") {
		resp.Header.Del("Content-Security-Policy")
		resp.Header.Del("X-Frame-Options")
	}
}

func injectInternalProxyAuth(req *http.Request, route repo.GatewayRoute) {
	if strings.HasPrefix(route.PathPattern, "/nodered/") || strings.HasPrefix(route.PathPattern, "/eventflow/") {
		if token := secrets.InternalToken("NODERED_INTERNAL_TOKEN"); token != "" {
			req.Header.Set(internalProxyTokenHeader, token)
		}
	}
}

// 独立打包的子应用（同源 iframe 承载），目录位于 WebDir/<name>，需要各自的 SPA index 回退
func serveSPAFile(w http.ResponseWriter, r *http.Request, filePath string) {
	if strings.HasSuffix(filePath, "index.html") {
		w.Header().Set("Cache-Control", "no-cache")
	} else if strings.Contains(filepath.ToSlash(filePath), "/assets/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	}
	http.ServeFile(w, r, filePath)
}

func (g *HTTPGateway) serveStatic(w http.ResponseWriter, r *http.Request) {
	if g.config.LocalFrontendDev && strings.TrimSpace(g.config.FrontendDevUrl) != "" {
		g.proxyFrontendDev(w, r)
		return
	}
	webDir := g.config.WebDir
	if webDir == "" {
		webDir = "./web"
	}
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}
	filePath := filepath.Join(webDir, filepath.FromSlash(path))
	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		serveSPAFile(w, r, filePath)
		return
	}
	indexPath := filepath.Join(webDir, "index.html")
	if _, err := os.Stat(indexPath); err == nil {
		serveSPAFile(w, r, indexPath)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte("<!doctype html><html><body><div id=\"root\">frontend build not found</div></body></html>"))
}

func (g *HTTPGateway) serveHealth(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code": 200,
		"msg":  "success",
		"data": map[string]any{
			"status": "ok",
			"name":   "backend",
		},
	})
}

func (g *HTTPGateway) proxyFrontendDev(w http.ResponseWriter, r *http.Request) {
	g.proxyDevUpstream(w, r, g.config.FrontendDevUrl)
}

// 开发态反向代理到本机前端 dev server（主应用/子应用共用）；
// WebSocket upgrade（Vite HMR）由 httputil.ReverseProxy 原生透传。
func (g *HTTPGateway) proxyDevUpstream(w http.ResponseWriter, r *http.Request, upstream string) {
	target, err := url.Parse(strings.TrimSpace(upstream))
	if err != nil || (target.Scheme != "http" && target.Scheme != "https") || target.Host == "" {
		http.Error(w, "invalid frontend dev upstream", http.StatusBadGateway)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = r.Host
		req.Header.Set("X-Forwarded-Host", r.Host)
		req.Header.Set("X-Forwarded-Proto", forwardedProto(r))
		// Vite dev server rejects oversized Cookie headers with 431; static/HMR assets do not need them.
		req.Header.Del("Cookie")
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}

func isInternalPath(path string) bool {
	return strings.HasPrefix(path, "/api/core/") ||
		strings.HasPrefix(path, "/openapi/v1/") ||
		path == "/healthz" ||
		path == "/readyz" ||
		path == "/metrics"
}

func routeMatches(route repo.GatewayRoute, method, path string) bool {
	if len(route.Methods) > 0 {
		ok := false
		for _, m := range route.Methods {
			if strings.EqualFold(m, method) || m == "*" {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	if route.PathPattern == path {
		return true
	}
	if strings.HasSuffix(route.PathPattern, "/**") {
		prefix := strings.TrimSuffix(route.PathPattern, "/**")
		return path == prefix || strings.HasPrefix(path, prefix+"/")
	}
	return false
}

func rewriteProxyPath(path string, route repo.GatewayRoute) string {
	if route.RewritePath != "" {
		return route.RewritePath
	}
	targetPath := strings.TrimSpace(route.TargetPath)
	if targetPath == "" {
		targetPath = "/"
	}
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = "/" + targetPath
	}
	if targetPath != "/" {
		targetPath = strings.TrimRight(targetPath, "/")
	}
	if !route.StripPrefix {
		if targetPath == "" {
			return path
		}
		return joinProxyPath(targetPath, path)
	}
	prefix := strings.TrimSuffix(route.PathPattern, "/**")
	rest := strings.TrimPrefix(path, prefix)
	if rest == "" {
		rest = "/"
	}
	if !strings.HasPrefix(rest, "/") {
		rest = "/" + rest
	}
	return joinProxyPath(targetPath, rest)
}

func joinProxyPath(base, rest string) string {
	if base == "" || base == "/" {
		return rest
	}
	if rest == "" || rest == "/" {
		return base + "/"
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(rest, "/")
}

func splitPathSegments(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

func isDeniedTarget(target *url.URL) bool {
	host := strings.Split(target.Host, ":")[0]
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate()
}

func forwardedProto(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func WithHTTPGateway(gateway *HTTPGateway) func(*http.Server) {
	return func(svr *http.Server) {
		if svr == nil || svr.Handler == nil || gateway == nil {
			return
		}
		next := svr.Handler
		svr.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if gateway.Handle(w, r) {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func WaitReady(ctx context.Context, app interface{ Ready(context.Context) error }, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if err := app.Ready(ctx); err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return errors.New("timeout waiting gateway dependencies")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}
