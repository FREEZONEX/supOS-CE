// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package system

import (
	"context"
	"os"
	"strings"

	respx "backend/internal/httpx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	unsTreeLoadModeAuto = "auto"
	unsTreeLoadModeLazy = "lazy"
	unsTreeLoadModeFull = "full"

	unsTreeLoadMode      = unsTreeLoadModeAuto
	unsTreeAutoThreshold = int64(2000)

	systemQualityField   = "_quality"
	systemTimestampField = "_timestamp"
	systemMQTTTCPPort    = "1883"
	systemMQTTWSPort     = "8083"
	systemMQTTWSSPort    = "8084"
)

type ConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigLogic {
	return &ConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigLogic) Config() (resp *types.Envelope, err error) {
	lazyTree, activeNodeCount := l.unsTreeLoadingConfig()
	productVersion := l.productVersion()
	return respx.Envelope(map[string]any{
		"appTitle":                 "Tier0",
		"authEnable":               true,
		"loginPath":                "/tier0-login",
		"productVersion":           productVersion,
		"version":                  productVersion,
		"lang":                     configuredLanguage(),
		"lazyTree":                 lazyTree,
		"unsTreeLoadMode":          unsTreeLoadMode,
		"unsTreeAutoThreshold":     unsTreeAutoThreshold,
		"unsTreeActiveNodeCount":   activeNodeCount,
		"enableAutoCategorization": true,
		"qualityName":              systemQualityField,
		"timestampName":            systemTimestampField,
		"mqttTcpPort":              configuredPort("OS_MQTT_TCP_PORT", systemMQTTTCPPort),
		"mqttWebsocketPort":        configuredPort("OS_MQTT_WEBSOCKET_PORT", systemMQTTWSPort),
		"mqttWebsocketTslPort":     configuredPort("OS_MQTT_WEBSOCKET_TSL_PORT", systemMQTTWSSPort),
		"containerMap":             systemContainerMap(),
	}), nil
}

func configuredPort(key, fallback string) string {
	if port := strings.TrimSpace(os.Getenv(key)); port != "" {
		return port
	}
	return fallback
}

func (l *ConfigLogic) productVersion() string {
	if l.svcCtx != nil {
		if version := strings.TrimSpace(l.svcCtx.Config.ProductVersion); version != "" {
			return version
		}
	}
	if version := strings.TrimSpace(os.Getenv("PRODUCT_VERSION")); version != "" {
		return version
	}
	return "dev"
}

func (l *ConfigLogic) unsTreeLoadingConfig() (bool, int64) {
	if l.svcCtx == nil || l.svcCtx.App == nil || l.svcCtx.App.UNS == nil {
		logx.WithContext(l.ctx).Error("UNS tree loading config fallback: UNS service is unavailable")
		return true, 0
	}
	count, err := l.svcCtx.App.UNS.ActiveNodeCount(l.ctx)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UNS tree loading config fallback: count active UNS nodes failed: %v", err)
		return true, 0
	}
	return resolveUnsTreeLazy(unsTreeLoadMode, count), count
}

func resolveUnsTreeLazy(mode string, activeNodeCount int64) bool {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case unsTreeLoadModeLazy:
		return true
	case unsTreeLoadModeFull:
		return false
	case unsTreeLoadModeAuto:
		return activeNodeCount > unsTreeAutoThreshold
	default:
		return true
	}
}

func configuredLanguage() string {
	for _, key := range []string{"LANGUAGE", "OS_LANG"} {
		if lang := normalizeLanguage(os.Getenv(key)); lang != "" {
			return lang
		}
	}
	return "en-US"
}

func normalizeLanguage(lang string) string {
	switch strings.ToLower(strings.ReplaceAll(strings.TrimSpace(lang), "_", "-")) {
	case "zh", "zh-cn":
		return "zh-CN"
	case "en", "en-us":
		return "en-US"
	default:
		return ""
	}
}

func systemContainerMap() map[string]any {
	return map[string]any{
		"sourceflow": map[string]any{
			"name":   "sourceflow",
			"envMap": map[string]any{"service_is_show": true, "service_redirect_url": "/nodered/home/"},
		},
		"eventflow": map[string]any{
			"name":   "eventflow",
			"envMap": map[string]any{"service_is_show": true, "service_redirect_url": "/eventflow/home/"},
		},
	}
}
