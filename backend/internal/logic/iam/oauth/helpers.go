package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	sysconfig "backend/internal/common/config"
	"backend/internal/common/constants"
	"backend/internal/repo/relationDB"
)

const (
	authorizationCodeTTL = 5 * time.Minute
	defaultOAuthScope    = "openid"
)

type OAuthError struct {
	Status      int
	Code        string
	Description string
}

func (e *OAuthError) Error() string {
	if e == nil {
		return ""
	}
	if strings.TrimSpace(e.Description) != "" {
		return e.Description
	}
	return e.Code
}

func newOAuthError(status int, code, description string) *OAuthError {
	return &OAuthError{
		Status:      status,
		Code:        strings.TrimSpace(code),
		Description: strings.TrimSpace(description),
	}
}

func loginRedirectURL(requestURI string) string {
	loginPath := sysconfig.NormalizeLoginPath(os.Getenv("SYS_OS_LOGIN_PATH"))
	target, err := url.Parse(loginPath)
	if err != nil {
		target = &url.URL{Path: sysconfig.DefaultLoginPath}
	}
	query := target.Query()
	query.Set("redirectUri", requestURI)
	target.RawQuery = query.Encode()
	return target.String()
}

func appendAuthorizeCode(redirectURI, code, state string) (string, error) {
	target, err := url.Parse(strings.TrimSpace(redirectURI))
	if err != nil {
		return "", err
	}
	query := target.Query()
	query.Set("code", code)
	if strings.TrimSpace(state) != "" {
		query.Set("state", strings.TrimSpace(state))
	}
	target.RawQuery = query.Encode()
	return target.String(), nil
}

func normalizeScope(scope string) string {
	scope = strings.TrimSpace(scope)
	if scope == "" {
		return defaultOAuthScope
	}
	return scope
}

func splitRedirectURIs(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == ';'
	})
	result := make([]string, 0, len(fields))
	for _, field := range fields {
		if trimmed := strings.TrimSpace(field); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func redirectURIMatches(client *relationDB.IamOAuthClient, redirectURI string) bool {
	if client == nil {
		return false
	}
	redirectURI = strings.TrimSpace(redirectURI)
	if redirectURI == "" {
		return false
	}
	for _, allowed := range splitRedirectURIs(client.RedirectURIs) {
		if strings.EqualFold(allowed, redirectURI) {
			return true
		}
	}
	return false
}

func accessTokenTTL() time.Duration {
	return time.Duration(constants.TokenMaxAge) * time.Second
}

func newOpaqueToken(size int) (string, error) {
	if size <= 0 {
		size = 32
	}
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func bearerToken(r *http.Request) string {
	if r == nil {
		return ""
	}
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(header) < len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}

func writeOAuthError(w http.ResponseWriter, err error) bool {
	if w == nil || err == nil {
		return false
	}
	var oauthErr *OAuthError
	if !errors.As(err, &oauthErr) || oauthErr == nil {
		return false
	}
	status := oauthErr.Status
	if status == 0 {
		status = http.StatusBadRequest
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + oauthErr.Code + `","error_description":"` + escapeJSON(oauthErr.Description) + `"}`))
	return true
}

func escapeJSON(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\r", `\r`, "\t", `\t`)
	return replacer.Replace(value)
}
