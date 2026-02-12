package model

import (
	"backend/share/app/util"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gitee.com/unitedrhino/share/errors"
	"gopkg.in/yaml.v3"
)

// NewFeatureModel 新功能模型
// 用于描述和管理应用的新功能特性
type NewFeatureModel struct {
	// Name 功能名称
	Name string `json:"name" yaml:"name"`

	// Description 功能描述
	Description string `json:"description" yaml:"description"`

	// ImagePath 镜像文件路径（本地文件路径fileId）
	ImagePath string `json:"imagePath" yaml:"imagePath"`

	// ImageUrl 镜像URL（远程镜像地址）
	ImageUrl string `json:"imageUrl" yaml:"imageUrl"`

	// ComposeYaml Docker Compose配置内容
	ComposeYaml string `json:"composeYaml" yaml:"composeYaml"`

	Menu *MenuModel `json:"menu" yaml:"menu"`

	InstallTime string `json:"installTime,omitempty"`
}

// Validate 验证模型数据的有效性
func (m *NewFeatureModel) Validate() error {
	if m.Name == "" {
		return errors.Parameter.WithMsg(fmt.Sprintf("app.name.empty"))
	}

	// 验证非法字符
	if !ValidateFilenameBasic(m.Name) {
		return fmt.Errorf("app.name.invalid")
	}
	// 验证该名称是否已经存在
	fullPath := filepath.Join(util.AppInstalledDir, m.Name)
	if _, err := os.Stat(fullPath); err == nil {
		return errors.Parameter.WithMsg(fmt.Sprintf("app.name.exist"))
	}

	// 镜像路径和URL至少需要提供一个
	if m.ImagePath == "" && m.ImageUrl == "" {
		return errors.Parameter.WithMsg(fmt.Sprintf("app.image.empty"))
	}

	// 验证菜单URL格式（如果提供了菜单URL）
	if m.Menu == nil {
		return errors.Parameter.WithMsg("menu.empty")
	}
	if err := m.Menu.Validate(); err != nil {
		return err
	}

	// 验证Docker Compose配置（如果提供了）
	if m.ComposeYaml != "" {
		if err := validateComposeYaml(m.ComposeYaml); err != nil {
			return err
		}
	}

	return nil
}

// validateComposeYaml 验证 Docker Compose YAML 配置
func validateComposeYaml(composeYaml string) error {
	if composeYaml == "" {
		return errors.Parameter.WithMsg("Docker Compose 配置不能为空")
	}

	// 1. 验证 YAML 格式
	var composeConfig map[string]interface{}
	if err := yaml.Unmarshal([]byte(composeYaml), &composeConfig); err != nil {
		return errors.Parameter.WithMsg(fmt.Sprintf("YAML 格式错误: %v", err))
	}

	// 2. 验证 services 部分
	services, ok := composeConfig["services"].(map[string]interface{})
	if !ok || len(services) == 0 {
		return errors.Parameter.WithMsg("必须包含至少一个 service")
	}

	// 3. 验证只能有一个服务
	if len(services) > 1 {
		return errors.Parameter.WithMsg("只能包含一个服务（container）")
	}

	// 4. 验证不能包含 network 声明
	if _, hasNetwork := composeConfig["networks"]; hasNetwork {
		return errors.Parameter.WithMsg("不能包含 network 声明，系统会自动管理网络")
	}

	// 5. 验证每个服务的配置
	for serviceName, serviceConfig := range services {
		serviceMap, ok := serviceConfig.(map[string]interface{})
		if !ok {
			return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 配置格式错误", serviceName))
		}

		// 验证必需的服务字段
		if _, ok := serviceMap["image"]; !ok {
			// 检查是否有 build 配置
			if _, hasBuild := serviceMap["build"]; !hasBuild {
				return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 必须指定 image 或 build 配置", serviceName))
			}
		}

		// 验证 container_name 配置
		if containerName, ok := serviceMap["container_name"].(string); ok && containerName != "" {
			// 验证 container_name 格式
			if err := validateContainerName(containerName); err != nil {
				return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 的 container_name 格式错误: %v", serviceName, err))
			}
		} else {
			// 如果没有 container_name，使用服务名称作为容器名称
			log.Printf("[app]: 服务 %s 没有指定 container_name，将使用服务名称作为容器名称", serviceName)
		}

		// 验证端口配置（可选）
		if portsConfig, ok := serviceMap["ports"].([]interface{}); ok {
			for i, portItem := range portsConfig {
				portStr, ok := portItem.(string)
				if !ok {
					return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 的端口配置 %d 格式错误", serviceName, i))
				}

				// 验证端口格式
				if err := validatePortFormat(portStr); err != nil {
					return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 的端口配置错误: %v", serviceName, err))
				}
			}
		}

		// 验证不能包含自定义网络配置
		if networksConfig, ok := serviceMap["networks"].([]interface{}); ok {
			if len(networksConfig) > 0 {
				return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 不能包含 networks 配置，系统会自动管理网络", serviceName))
			}
		}
		if networksConfig, ok := serviceMap["networks"].(map[string]interface{}); ok {
			if len(networksConfig) > 0 {
				return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 不能包含 networks 配置，系统会自动管理网络", serviceName))
			}
		}

		// 验证其他重要字段
		validateServiceOptionalFields(serviceName, serviceMap)
	}

	return nil
}

// validateContainerName 验证容器名称格式
func validateContainerName(containerName string) error {
	if containerName == "" {
		return fmt.Errorf("容器名称不能为空")
	}

	// 检查长度
	if len(containerName) > 255 {
		return fmt.Errorf("容器名称长度不能超过255个字符")
	}

	// 检查非法字符（Docker 容器名称限制）
	// Docker 容器名称只能包含: [a-zA-Z0-9][a-zA-Z0-9_.-]
	for i, ch := range containerName {
		if i == 0 {
			// 第一个字符必须是字母或数字
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
				return fmt.Errorf("容器名称必须以字母或数字开头")
			}
		} else {
			// 后续字符可以是字母、数字、下划线、点或短横线
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') ||
				ch == '_' || ch == '.' || ch == '-') {
				return fmt.Errorf("容器名称包含非法字符 '%c'，只能包含字母、数字、下划线、点和短横线", ch)
			}
		}
	}

	// 检查保留名称
	reservedNames := []string{"none", "host", "bridge"}
	for _, reserved := range reservedNames {
		if strings.EqualFold(containerName, reserved) {
			return fmt.Errorf("容器名称不能使用保留名称 '%s'", reserved)
		}
	}

	return nil
}

// validatePortFormat 验证端口格式
func validatePortFormat(portStr string) error {
	portStr = strings.TrimSpace(portStr)
	if portStr == "" {
		return fmt.Errorf("端口字符串为空")
	}

	// 检查格式: 可能的形式有:
	// 1. "8080" - 仅容器端口
	// 2. "8080:80" - 主机端口:容器端口
	// 3. "127.0.0.1:8080:80" - IP:主机端口:容器端口
	// 4. "8080:80/tcp" - 带协议
	// 5. "127.0.0.1:8080:80/udp" - 完整格式

	// 分割协议部分
	parts := strings.Split(portStr, "/")
	if len(parts) > 2 {
		return fmt.Errorf("端口格式错误: %s", portStr)
	}

	portWithoutProtocol := parts[0]
	if len(parts) == 2 {
		protocol := strings.ToLower(parts[1])
		if protocol != "tcp" && protocol != "udp" {
			return fmt.Errorf("不支持的协议: %s", protocol)
		}
	}

	// 分割端口部分
	portParts := strings.Split(portWithoutProtocol, ":")
	if len(portParts) > 3 {
		return fmt.Errorf("端口格式错误: %s", portStr)
	}

	// 验证每个端口号
	for _, portPart := range portParts {
		// 如果是 IP 地址，跳过验证
		if isIPAddress(portPart) {
			continue
		}

		// 验证端口号
		port, err := strconv.Atoi(portPart)
		if err != nil {
			return fmt.Errorf("无效的端口号: %s", portPart)
		}

		if port < 1 || port > 65535 {
			return fmt.Errorf("端口号超出范围 (1-65535): %d", port)
		}
	}

	return nil
}

// isIPAddress 检查字符串是否是 IP 地址
func isIPAddress(str string) bool {
	// 简单检查是否是 IPv4 地址格式
	parts := strings.Split(str, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// validateServiceOptionalFields 验证服务的可选字段
func validateServiceOptionalFields(serviceName string, serviceConfig map[string]interface{}) {
	// 验证环境变量
	if env, ok := serviceConfig["environment"]; ok {
		switch envVal := env.(type) {
		case map[string]interface{}:
			// 格式正确
		case []interface{}:
			// 数组格式
			for i, item := range envVal {
				if _, ok := item.(string); !ok {
					log.Printf("[app]: 警告: 服务 %s 的环境变量 %d 格式错误", serviceName, i)
				}
			}
		default:
			log.Printf("[app]: 警告: 服务 %s 的环境变量格式错误", serviceName)
		}
	}

	// 验证卷挂载
	if volumes, ok := serviceConfig["volumes"]; ok {
		if volumeList, ok := volumes.([]interface{}); ok {
			for i, volume := range volumeList {
				if _, ok := volume.(string); !ok {
					log.Printf("[app]: 警告: 服务 %s 的卷配置 %d 格式错误", serviceName, i)
				}
			}
		} else {
			log.Printf("[app]: 警告: 服务 %s 的卷配置格式错误", serviceName)
		}
	}

	// 验证重启策略
	if restart, ok := serviceConfig["restart"].(string); ok {
		validRestartPolicies := []string{"no", "always", "on-failure", "unless-stopped"}
		valid := false
		for _, policy := range validRestartPolicies {
			if restart == policy {
				valid = true
				break
			}
		}
		if !valid {
			log.Printf("[app]: 警告: 服务 %s 的重启策略 '%s' 无效", serviceName, restart)
		}
	}
}

func ValidateFilenameBasic(filename string) bool {
	if filename == "" {
		return false
	}

	// 检查长度
	if len(filename) > 255 {
		return false
	}

	// 检查是否包含路径分隔符
	if strings.ContainsAny(filename, `/\:*?"<>|`) {
		return false
	}

	// 检查首尾字符
	if strings.HasPrefix(filename, ".") || strings.HasSuffix(filename, ".") {
		return false
	}

	// 检查空格
	if strings.HasPrefix(filename, " ") || strings.HasSuffix(filename, " ") {
		return false
	}

	return true
}

// GetImageSource 获取镜像来源
// 返回: "local" - 本地文件, "remote" - 远程URL, "" - 无镜像
func (m *NewFeatureModel) GetImageSource() string {
	if m.ImagePath != "" {
		return "local"
	} else if m.ImageUrl != "" {
		return "remote"
	}
	return ""
}

// GetImageIdentifier 获取镜像标识符
// 优先使用本地路径，其次使用远程URL
func (m *NewFeatureModel) GetImageIdentifier() string {
	if m.ImagePath != "" {
		return m.ImagePath
	}
	return m.ImageUrl
}

// LoadFromFile 从YAML文件加载模型
func LoadFromFile(filePath string) (*NewFeatureModel, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.Parameter.WithMsg(fmt.Sprintf("配置文件不存在: %s", filePath))
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Parameter.WithMsg(fmt.Sprintf("读取配置文件失败: %v", err))
	}

	// 解析YAML
	var model NewFeatureModel
	if err := yaml.Unmarshal(data, &model); err != nil {
		return nil, errors.Parameter.WithMsg(fmt.Sprintf("解析YAML配置失败: %v", err))
	}

	return &model, nil
}

// SaveToFile 将模型保存到YAML文件
func (m *NewFeatureModel) SaveToFile(filePath string) error {

	// 转换为YAML
	data, err := yaml.Marshal(m)
	if err != nil {
		return errors.Parameter.WithMsg(fmt.Sprintf("序列化YAML失败: %v", err))
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Parameter.WithMsg(fmt.Sprintf("创建目录失败: %v", err))
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return errors.Parameter.WithMsg(fmt.Sprintf("写入文件失败: %v", err))
	}

	return nil
}

// ToAppConfig 转换为AppConfig（如果需要与其他系统集成）
func (m *NewFeatureModel) ToAppConfig() *AppConfig {
	return &AppConfig{
		Name:        m.Name,
		Description: m.Description,
		// 其他字段根据需要进行映射
	}
}
