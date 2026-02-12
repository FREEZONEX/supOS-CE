package app

import (
	"backend/share/app/adapter"
	"backend/share/app/model"
	"backend/share/app/util"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"gopkg.in/yaml.v3"
)

//
/**
 * TODO 1. iconUrl、description字段非必填校验
*       2. 异常信息国际化
*       3. 容器名、端口占用校验
*/

func InstallFeature(fc *model.NewFeatureModel) error {
	if err := fc.Validate(); err != nil {
		return err
	}

	// 加载镜像
	switch fc.GetImageSource() {
	case "local":
		{
			if err := adapter.LoadImageFromLocal(fc.ImagePath); err != nil {
				return err
			}
		}
	case "remote":
		{
			if err := adapter.PullImageFromRemote(fc.ImageUrl); err != nil {
				return err
			}
		}
	}
	var containerName string
	var containerPorts []int

	if fc.ComposeYaml != "" {
		// 从 Compose YAML 中提取容器名称和端口
		name, ports, err := extractContainerInfoFromCompose(fc.ComposeYaml)
		if err != nil {
			log.Printf("[app]: 警告: 从 Compose YAML 提取容器信息失败: %v", err)
			// 不因为提取失败而停止安装
		} else {
			containerName = name
			containerPorts = ports
			log.Printf("[app]: 从 Compose YAML 提取到容器信息: 名称=%s, 端口=%v", containerName, containerPorts)
		}
		// 部署docker compose文件
		defaultNetwork := "tier0_edge_network"
		if err := adapter.DeployComposeYaml(fc.ComposeYaml, defaultNetwork); err != nil {
			return err
		}
	}
	var port int
	if len(containerPorts) > 0 {
		port = containerPorts[0]
	} else {
		port = 8080
	}
	// 注册菜单
	if err := adapter.SaveMenu(fc.Menu, containerName, port); err != nil {
		return err
	}
	// 将应用配置写到本地文件,文件名为fc.Name.json
	return saveFeatureConfig(fc)
}

func extractContainerInfoFromCompose(composeYaml string) (string, []int, error) {

	// 解析 YAML
	var composeConfig map[string]interface{}
	if err := yaml.Unmarshal([]byte(composeYaml), &composeConfig); err != nil {
		return "", nil, errors.Parameter.WithMsg(fmt.Sprintf("解析 Compose YAML 失败: %v", err))
	}

	// 检查 services 部分
	services, ok := composeConfig["services"].(map[string]interface{})
	if !ok || len(services) == 0 {
		return "", nil, errors.Parameter.WithMsg(fmt.Sprintf("Compose YAML 中没有找到 services 部分"))
	}

	// 获取第一个服务（通常只有一个主要服务）
	var firstServiceName string
	var firstServiceConfig interface{}
	for name, config := range services {
		firstServiceName = name
		firstServiceConfig = config
		break
	}

	// 解析服务配置
	serviceConfig, ok := firstServiceConfig.(map[string]interface{})
	if !ok {
		return "", nil, errors.Parameter.WithMsg(fmt.Sprintf("服务配置格式错误"))
	}

	// 提取容器名称
	containerName := ""
	if name, ok := serviceConfig["container_name"].(string); ok && name != "" {
		containerName = name
	} else {
		// 如果没有 container_name，使用服务名称
		containerName = firstServiceName
	}

	// 提取端口
	var ports []int
	if portsConfig, ok := serviceConfig["ports"].([]interface{}); ok {
		for _, portItem := range portsConfig {
			portStr, ok := portItem.(string)
			if !ok {
				continue
			}

			// 解析端口映射，格式可能是 "8080:80" 或 "8080"
			portParts := strings.Split(portStr, ":")
			if len(portParts) == 0 {
				continue
			}

			// 获取容器内部端口（通常是第二部分）
			var containerPortStr string
			if len(portParts) == 2 {
				// 格式: "主机端口:容器端口"
				containerPortStr = portParts[1]
			} else {
				// 格式: "容器端口" 或 "主机端口"
				containerPortStr = portParts[0]
			}

			// 移除可能的协议后缀（如 /tcp, /udp）
			if slashIndex := strings.Index(containerPortStr, "/"); slashIndex != -1 {
				containerPortStr = containerPortStr[:slashIndex]
			}

			// 转换为整数
			port, err := parsePort(containerPortStr)
			if err != nil {
				log.Printf("[app]: 警告: 无法解析端口 %s: %v", containerPortStr, err)
				continue
			}

			ports = append(ports, port)
		}
	}

	// 如果没有找到 ports 配置，尝试从 expose 中获取
	if len(ports) == 0 {
		if exposeConfig, ok := serviceConfig["expose"].([]interface{}); ok {
			for _, exposeItem := range exposeConfig {
				portStr, ok := exposeItem.(string)
				if !ok {
					continue
				}

				port, err := parsePort(portStr)
				if err != nil {
					log.Printf("[app]: 警告: 无法解析 expose 端口 %s: %v", portStr, err)
					continue
				}

				ports = append(ports, port)
			}
		}
	}

	log.Printf("[app]: 从 Compose YAML 提取信息: 容器名称=%s, 端口=%v", containerName, ports)
	return containerName, ports, nil
}

// parsePort 解析端口字符串
func parsePort(portStr string) (int, error) {
	portStr = strings.TrimSpace(portStr)
	if portStr == "" {
		return 0, errors.Parameter.WithMsg(fmt.Sprintf("端口字符串为空"))
	}

	// 尝试转换为整数
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, errors.Parameter.WithMsg(fmt.Sprintf("无效的端口号: %s", portStr))
	}

	// 验证端口范围
	if port < 1 || port > 65535 {
		return 0, errors.Parameter.WithMsg(fmt.Sprintf("端口号超出范围 (1-65535): %d", port))
	}

	return port, nil
}

// saveFeatureConfig 将功能配置保存到本地 JSON 文件
func saveFeatureConfig(fc *model.NewFeatureModel) error {

	// 1. 确保配置目录存在
	if err := os.MkdirAll(util.AppInstalledDir, 0755); err != nil {
		return errors.Parameter.WithMsg(fmt.Sprintf("创建配置目录失败: %v", err))
	}

	// 2. 构建文件路径
	// 使用功能名称作为文件名，确保文件名安全
	safeFileName := makeFileNameSafe(fc.Name)
	configPath := filepath.Join(util.AppInstalledDir, safeFileName+".json")

	// 设置安装时间为当前时间
	fc.InstallTime = time.Now().Format("2006-01-02 15:04:05")
	// 3. 转换为 JSON
	configJSON, err := json.MarshalIndent(fc, "", "  ")
	if err != nil {
		return errors.Parameter.WithMsg(fmt.Sprintf("序列化配置失败: %v", err))
	}

	// 4. 写入文件
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return errors.Parameter.WithMsg(fmt.Sprintf("写入配置文件失败: %v", err))
	}

	// 5. 记录日志
	log.Printf("[app]: 功能配置已保存: %s, 安装时间: %s", fc.Name, fc.InstallTime)
	return nil
}

func makeFileNameSafe(fileName string) string {
	// 定义非法字符
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	safeName := fileName

	// 替换非法字符为下划线
	for _, char := range invalidChars {
		safeName = strings.ReplaceAll(safeName, char, "_")
	}

	// 移除首尾空格
	safeName = strings.TrimSpace(safeName)

	// 如果文件名为空，使用默认名称
	if safeName == "" {
		safeName = "unnamed-feature"
	}

	// 限制文件名长度
	if len(safeName) > 100 {
		safeName = safeName[:100]
	}

	return safeName
}

// ValidateComposeYaml 验证 Docker Compose YAML 格式和端口占用
func ValidateComposeYaml(composeYaml string) error {
	if composeYaml == "" {
		return errors.Parameter.WithMsg("Docker Compose 配置不能为空")
	}

	// 1. 验证 YAML 格式
	var composeConfig map[string]interface{}
	if err := yaml.Unmarshal([]byte(composeYaml), &composeConfig); err != nil {
		return errors.Parameter.WithMsg(fmt.Sprintf("YAML 格式错误: %v", err))
	}

	// 2. 验证必需字段
	if _, ok := composeConfig["version"]; !ok {
		log.Printf("[app]: 警告: Compose YAML 中没有指定 version 字段")
		// 不强制要求 version，但记录警告
	}

	// 3. 验证 services 部分
	services, ok := composeConfig["services"].(map[string]interface{})
	if !ok || len(services) == 0 {
		return errors.Parameter.WithMsg("必须包含至少一个 service")
	}

	// 4. 验证每个服务的配置
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

		// 验证端口配置
		if err := validateServicePorts(serviceName, serviceMap); err != nil {
			return err
		}

		// 验证其他重要字段（可选）
		validateOptionalFields(serviceName, serviceMap)
	}

	// 5. 验证网络配置（如果存在）
	if networks, ok := composeConfig["networks"].(map[string]interface{}); ok {
		for networkName, networkConfig := range networks {
			if err := validateNetworkConfig(networkName, networkConfig); err != nil {
				log.Printf("[app]: 警告: 网络 %s 配置验证失败: %v", networkName, err)
			}
		}
	}

	// 6. 验证卷配置（如果存在）
	if volumes, ok := composeConfig["volumes"].(map[string]interface{}); ok {
		for volumeName, volumeConfig := range volumes {
			if err := validateVolumeConfig(volumeName, volumeConfig); err != nil {
				log.Printf("[app]: 警告: 卷 %s 配置验证失败: %v", volumeName, err)
			}
		}
	}

	log.Printf("[app]: Docker Compose YAML 格式验证通过")
	return nil
}

// validateServicePorts 验证服务端口配置
func validateServicePorts(serviceName string, serviceConfig map[string]interface{}) error {
	// 检查端口配置
	if portsConfig, ok := serviceConfig["ports"].([]interface{}); ok {
		for i, portItem := range portsConfig {
			portStr, ok := portItem.(string)
			if !ok {
				return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 的端口配置 %d 格式错误", serviceName, i))
			}

			// 验证端口格式
			if err := validatePortFormat(portStr); err != nil {
				return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 的端口配置错误: %v", serviceName, err))
			}

			// 检查端口是否被占用
			if err := checkPortAvailability(portStr); err != nil {
				return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 的端口 %s 已被占用", serviceName, portStr))
			}
		}
	}

	// 检查 expose 配置
	if exposeConfig, ok := serviceConfig["expose"].([]interface{}); ok {
		for i, exposeItem := range exposeConfig {
			portStr, ok := exposeItem.(string)
			if !ok {
				return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 的 expose 配置 %d 格式错误", serviceName, i))
			}

			// 验证端口格式
			if err := validatePortFormat(portStr); err != nil {
				return errors.Parameter.WithMsg(fmt.Sprintf("服务 %s 的 expose 端口错误: %v", serviceName, err))
			}
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
	// 检查 IPv4
	if ip := net.ParseIP(str); ip != nil {
		return true
	}

	// 检查是否是主机名（包含字母）
	for _, r := range str {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return false
		}
	}

	return false
}

// checkPortAvailability 检查端口是否可用
func checkPortAvailability(portStr string) error {
	// 解析端口字符串，获取主机端口
	portParts := strings.Split(strings.Split(portStr, "/")[0], ":")

	var hostPortStr string
	if len(portParts) == 1 {
		// 格式: "容器端口" - 不检查主机端口占用
		return nil
	} else if len(portParts) == 2 {
		// 格式: "主机端口:容器端口"
		hostPortStr = portParts[0]
	} else if len(portParts) == 3 {
		// 格式: "IP:主机端口:容器端口"
		hostPortStr = portParts[1]
	}

	// 如果是 IP 绑定，检查是否是本地 IP
	if len(portParts) == 3 {
		ip := portParts[0]
		if ip != "127.0.0.1" && ip != "0.0.0.0" && ip != "::" && ip != "localhost" {
			// 非本地 IP 绑定，不检查端口占用
			return nil
		}
	}

	// 解析主机端口
	hostPort, err := strconv.Atoi(hostPortStr)
	if err != nil {
		return fmt.Errorf("无效的主机端口: %s", hostPortStr)
	}

	// 检查端口是否被占用
	if isPortInUse(hostPort) {
		return fmt.Errorf("端口 %d 已被占用", hostPort)
	}

	return nil
}

// isPortInUse 检查端口是否被占用
func isPortInUse(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// 端口可能被占用或其他错误
		log.Printf("[app]: 端口 %d 检查: %v", port, err)
		return true
	}
	defer listener.Close()
	return false
}

// validateOptionalFields 验证可选字段
func validateOptionalFields(serviceName string, serviceConfig map[string]interface{}) {
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

// validateNetworkConfig 验证网络配置
func validateNetworkConfig(networkName string, networkConfig interface{}) error {
	// 网络配置可以是简单字符串或对象
	switch config := networkConfig.(type) {
	case string:
		// 简单字符串格式
		if config == "" {
			return fmt.Errorf("网络 %s 配置为空", networkName)
		}
	case map[string]interface{}:
		// 对象格式，验证可选字段
		if driver, ok := config["driver"].(string); ok {
			validDrivers := []string{"bridge", "overlay", "host", "none"}
			valid := false
			for _, d := range validDrivers {
				if driver == d {
					valid = true
					break
				}
			}
			if !valid {
				log.Printf("[app]: 警告: 网络 %s 的驱动 '%s' 可能不受支持", networkName, driver)
			}
		}
	default:
		return fmt.Errorf("网络 %s 配置格式错误", networkName)
	}

	return nil
}

// validateVolumeConfig 验证卷配置
func validateVolumeConfig(volumeName string, volumeConfig interface{}) error {
	// 卷配置可以是简单字符串或对象
	switch config := volumeConfig.(type) {
	case string:
		// 简单字符串格式
		if config == "" {
			return fmt.Errorf("卷 %s 配置为空", volumeName)
		}
	case map[string]interface{}:
		// 对象格式，验证可选字段
		if driver, ok := config["driver"].(string); ok {
			if driver == "" {
				log.Printf("[app]: 警告: 卷 %s 的驱动名为空", volumeName)
			}
		}
		// 可以添加更多验证逻辑
	default:
		return fmt.Errorf("卷 %s 配置格式错误", volumeName)
	}

	return nil
}
