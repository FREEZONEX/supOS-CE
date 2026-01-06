package system

import (
	"backend/share/spring"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	I18nUtils "backend/internal/common/I18nUtils"
	sysconfig "backend/internal/common/config"
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/common/utils/fileutil"
	"backend/internal/common/utils/runtimeutil"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v3"
)

type SystemConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取系统配置
func NewSystemConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SystemConfigLogic {
	return &SystemConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func init() {
	spring.RegisterLazy[*sysconfig.SystemConfig](func() *sysconfig.SystemConfig {
		logger := logx.WithContext(context.Background())
		cfg, buildErr := buildSystemConfig(logger)
		if buildErr != nil {
			// Log the error but still attempt to return whatever we have.
			logger.Errorf("failed to build system config: %v", buildErr)
		}
		return cfg
	})
}
func (l *SystemConfigLogic) SystemConfig() (resp *types.SystemConfigResp, err error) {
	cfg := spring.GetBean[*sysconfig.SystemConfig]()
	resp = mapSystemConfigToResp(l.ctx, cfg)
	return
}

const (
	systemServicePrefix   = "service_"
	activeServicesFile    = "active-services.txt"
	defaultComposePattern = "docker-compose"
)

func buildSystemConfig(logger logx.Logger) (*sysconfig.SystemConfig, error) {
	cfg := sysconfig.NewSystemConfig()

	cfg.AppTitle = firstNonEmpty(os.Getenv("SYS_OS_APP_TITLE"), cfg.AppTitle)
	cfg.Version = firstNonEmpty(os.Getenv("SYS_OS_VERSION"), firstNonEmpty(cfg.Version, constants.OSVersion))
	cfg.Lang = firstNonEmpty(os.Getenv("SYS_OS_LANG"), cfg.Lang)
	cfg.AuthEnable = boolEnv("SYS_OS_AUTH_ENABLE", cfg.AuthEnable)
	cfg.LLMType = firstNonEmpty(os.Getenv("SYS_OS_LLM_TYPE"), cfg.LLMType)
	cfg.PlatformType = firstNonEmpty(os.Getenv("SYS_OS_PLATFORM_TYPE"), cfg.PlatformType)
	cfg.EntranceURL = strings.TrimSpace(os.Getenv("SYS_OS_ENTRANCE_URL"))
	cfg.LoginPath = firstNonEmpty(os.Getenv("SYS_OS_LOGIN_PATH"), cfg.LoginPath)
	cfg.MQTTTCPPort = intEnv("SYS_OS_MQTT_TCP_PORT", cfg.MQTTTCPPort)
	cfg.MQTTWebsocketTSLPort = intEnv("SYS_OS_MQTT_WEBSOCKET_TSL_PORT", cfg.MQTTWebsocketTSLPort)
	cfg.MultipleTopic = boolEnv("SYS_OS_MULTIPLE_TOPIC", cfg.MultipleTopic)
	cfg.UseAliasPathAsTopic = boolEnv("SYS_OS_USE_ALIAS_PATH_AS_TOPIC", constants.UseAliasAsTopic)
	cfg.QualityName = firstNonEmpty(os.Getenv("SYS_OS_QUALITY_NAME"), firstNonEmpty(cfg.QualityName, constants.QosField))
	cfg.TimestampName = firstNonEmpty(os.Getenv("SYS_OS_TIMESTAMP_NAME"), firstNonEmpty(cfg.TimestampName, constants.SysFieldCreateTime))
	cfg.LazyTree = boolEnv("SYS_OS_LAZY_TREE", cfg.LazyTree)
	cfg.LDAPEnable = boolEnv("SYS_OS_LDAP_ENABLE", cfg.LDAPEnable)
	cfg.EnableAutoCategorization = boolEnv("SYS_OS_ENABLE_AUTO_CATEGORIZATION", cfg.EnableAutoCategorization)
	containerMap, err := loadContainerMap(logger)
	if err != nil {
		return cfg, err
	}
	cfg.ContainerMap = containerMap

	return cfg, nil
}

func mapSystemConfigToResp(ctx context.Context, cfg *sysconfig.SystemConfig) *types.SystemConfigResp {
	if cfg == nil {
		return &types.SystemConfigResp{}
	}
	resp := &types.SystemConfigResp{
		AppTitle:             cfg.AppTitle,
		Version:              firstNonEmpty(cfg.Version, constants.OSVersion),
		Lang:                 cfg.Lang,
		AuthEnable:           cfg.AuthEnable,
		LlmType:              cfg.LLMType,
		EnableAi:             len(strings.TrimSpace(os.Getenv("LLM_API_KEY"))) > 0,
		MqttTcpPort:          int64(cfg.MQTTTCPPort),
		MqttWebsocketTslPort: int64(cfg.MQTTWebsocketTSLPort),
		LoginPath:            cfg.LoginPath,
		PlatformType:         cfg.PlatformType,
		EntranceUrl:          cfg.EntranceURL,
		MultipleTopic:        cfg.MultipleTopic,
		UseAliasPathAsTopic:  cfg.UseAliasPathAsTopic,
		QualityName:          cfg.QualityName,
		TimestampName:        cfg.TimestampName,
		LazyTree:             cfg.LazyTree,
		LdapEnable:           cfg.LDAPEnable,
		AutoCategoryEnable:   cfg.EnableAutoCategorization,
	}

	if len(cfg.ContainerMap) > 0 {
		resp.ContainerMap = make(map[string]types.ContainerInfo, len(cfg.ContainerMap))
		for key, info := range cfg.ContainerMap {
			if info == nil {
				continue
			}
			envMap := make(map[string]interface{}, len(info.EnvMap))
			for k, v := range info.EnvMap {
				envMap[k] = v
			}

			description := strings.TrimSpace(info.Description)
			if description != "" && !strings.EqualFold(description, "null") {
				description = I18nUtils.GetMessageWithCtx(ctx, description)
			} else {
				description = ""
			}

			resp.ContainerMap[key] = types.ContainerInfo{
				Name:        info.Name,
				Version:     info.Version,
				Description: description,
				EnvMap:      envMap,
			}
		}
	}

	return resp
}

func loadContainerMap(logger logx.Logger) (map[string]*sysconfig.ContainerInfo, error) {
	containerMap, err := buildContainerMap()
	if err != nil {
		logger.Errorf("failed to load system container map: %v", err)
		return map[string]*sysconfig.ContainerInfo{}, err
	}
	return containerMap, nil
}

type composeFile struct {
	Services map[string]composeService `yaml:"services"`
}

type composeService struct {
	ContainerName string      `yaml:"container_name"`
	Image         string      `yaml:"image"`
	Environment   interface{} `yaml:"environment"`
}

func buildContainerMap() (map[string]*sysconfig.ContainerInfo, error) {
	composeData, err := loadComposeFile()
	if err != nil {
		return nil, fmt.Errorf("load compose file: %w", err)
	}
	if len(composeData) == 0 {
		return map[string]*sysconfig.ContainerInfo{}, nil
	}

	var compose composeFile
	if err := yaml.Unmarshal(composeData, &compose); err != nil {
		return nil, fmt.Errorf("unmarshal compose yaml: %w", err)
	}

	activeLine, err := loadActiveServicesLine()
	if err != nil {
		return nil, fmt.Errorf("load active services: %w", err)
	}

	result := make(map[string]*sysconfig.ContainerInfo)
	for serviceName, service := range compose.Services {
		containerName := strings.TrimSpace(service.ContainerName)
		if containerName == "" {
			containerName = serviceName
		}
		if containerName == "" {
			continue
		}
		if !isServiceActive(activeLine, containerName) {
			continue
		}

		envMap := filterServiceEnv(service.Environment)
		envMap[enums.ContainerEnvServiceIsShow.Name] = true

		version := extractVersion(service.Image)
		description := toString(envMap[enums.ContainerEnvServiceDescription.Name])

		result[containerName] = &sysconfig.ContainerInfo{
			Name:        containerName,
			Version:     version,
			Description: description,
			EnvMap:      envMap,
		}

		if strings.EqualFold(serviceName, "emqx") && isServiceActive(activeLine, "gmqtt") {
			result["gmqtt"] = &sysconfig.ContainerInfo{
				Name:        "gmqtt",
				Version:     version,
				Description: description,
				EnvMap:      cloneEnvMap(envMap),
			}
		}
	}

	return result, nil
}

func loadComposeFile() ([]byte, error) {
	candidates := []string{
		// strings.TrimSpace(os.Getenv("SYS_OS_COMPOSE_PATH")),
		"/app/go-edge/system",
	}
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		data, err := readComposeCandidate(candidate)
		if err != nil {
			return nil, err
		}
		if len(data) > 0 {
			return data, nil
		}
	}

	if runtimeutil.IsLocalEnv() {
		candidates := []string{
			filepath.Join("backend", "templates", "docker-compose-8c16g.yml"),
			filepath.Join("templates", "docker-compose-8c16g.yml"),
			filepath.Join("deploy", "docker", "run-env", "docker-compose.yml"),
			filepath.Join("deploy", "docker", "docker-compose.yml"),
		}
		for _, path := range candidates {
			if data, err := os.ReadFile(path); err == nil {
				return data, nil
			}
		}
		return nil, nil
	}

	systemDir := filepath.Join(fileutil.GetFileRootPath(), strings.Trim(constants.SystemRoot, "/"))
	return readComposeCandidate(systemDir)
}

func loadActiveServicesLine() (string, error) {
	if custom := strings.TrimSpace("/app/go-edge/system/active-services.txt"); custom != "" {
		data, err := os.ReadFile(custom)
		if err != nil {
			return "", err
		}
		return firstLine(data), nil
	}

	if runtimeutil.IsLocalEnv() {
		candidates := []string{
			filepath.Join("backend", "templates", activeServicesFile),
			filepath.Join("templates", activeServicesFile),
			filepath.Join("deploy", "docker", "run-env", "global", activeServicesFile),
			filepath.Join("deploy", "docker", "global", activeServicesFile),
		}
		for _, path := range candidates {
			if data, err := os.ReadFile(path); err == nil {
				return firstLine(data), nil
			}
		}
		return "", nil
	}

	systemDir := filepath.Join(fileutil.GetFileRootPath(), strings.Trim(constants.SystemRoot, "/"))
	path := filepath.Join(systemDir, activeServicesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return firstLine(data), nil
}

func readComposeCandidate(target string) ([]byte, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return nil, nil
	}

	info, err := os.Stat(target)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	if info.IsDir() {
		entries, err := os.ReadDir(target)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if strings.HasPrefix(name, defaultComposePattern) {
				return os.ReadFile(filepath.Join(target, name))
			}
		}
		return nil, nil
	}
	return os.ReadFile(target)
}

func firstLine(data []byte) string {
	content := strings.TrimSpace(string(data))
	if content == "" {
		return ""
	}
	for _, sep := range []string{"\r\n", "\n", "\r"} {
		if idx := strings.Index(content, sep); idx >= 0 {
			return strings.TrimSpace(content[:idx])
		}
	}
	return content
}

func isServiceActive(activeLine, containerName string) bool {
	if strings.TrimSpace(containerName) == "" {
		return false
	}
	if strings.TrimSpace(activeLine) == "" {
		return true
	}
	return strings.Contains(activeLine, containerName)
}

func filterServiceEnv(raw interface{}) map[string]any {
	rawMap := toEnvMap(raw)
	if len(rawMap) == 0 {
		return make(map[string]any)
	}
	result := make(map[string]any)
	for key, value := range rawMap {
		if strings.HasPrefix(key, systemServicePrefix) {
			result[key] = value
		}
	}
	return result
}

func toEnvMap(raw interface{}) map[string]any {
	result := make(map[string]any)
	switch v := raw.(type) {
	case map[string]interface{}:
		for key, value := range v {
			result[key] = value
		}
	case map[interface{}]interface{}:
		for key, value := range v {
			if keyStr, ok := key.(string); ok {
				result[keyStr] = value
			}
		}
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok {
				if parts := strings.SplitN(str, "=", 2); len(parts) == 2 {
					result[parts[0]] = parts[1]
				}
			}
		}
	case []string:
		for _, item := range v {
			if parts := strings.SplitN(item, "=", 2); len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}
	}
	return result
}

func cloneEnvMap(src map[string]any) map[string]any {
	if len(src) == 0 {
		return make(map[string]any)
	}
	clone := make(map[string]any, len(src))
	for k, v := range src {
		clone[k] = v
	}
	return clone
}

func extractVersion(image string) string {
	image = strings.TrimSpace(image)
	if image == "" {
		return ""
	}
	if idx := strings.LastIndex(image, ":"); idx >= 0 && idx < len(image)-1 {
		return image[idx+1:]
	}
	return ""
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		if value == nil {
			return ""
		}
		return fmt.Sprint(value)
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func boolEnv(key string, def bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.ParseBool(strings.TrimSpace(val)); err == nil {
			return parsed
		}
	}
	return def
}

func intEnv(key string, def int) int {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
			return parsed
		}
	}
	return def
}
