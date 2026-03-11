package model

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ComposeConfig 表示compose.yaml文件的完整结构
type ComposeConfig struct {
	APIVersion string        `yaml:"apiVersion"`
	Service    ServiceConfig `yaml:"service"`
}

// ServiceConfig 表示service配置
type ServiceConfig struct {
	MultiReplicas bool         `yaml:"multiReplicas"`
	Replicas      int          `yaml:"replicas"`
	Proxy         *ProxyConfig `yaml:"proxy,omitempty"`
	Ports         []PortConfig `yaml:"ports"`
	Containers    []Container  `yaml:"containers"`
}

// ProxyConfig 表示代理配置
type ProxyConfig struct {
	Paths []ProxyPath `yaml:"paths"`
}

// ProxyPath 表示代理路径配置
type ProxyPath struct {
	Path        string `yaml:"path"`
	ServicePort int    `yaml:"servicePort"`
}

// PortConfig 表示端口配置
type PortConfig struct {
	Name       string `yaml:"name"`
	Port       int    `yaml:"port"`
	Protocol   string `yaml:"protocol"`
	TargetPort int    `yaml:"targetPort"`
}

// Container 表示容器配置
type Container struct {
	Name         string          `yaml:"name"`
	Image        string          `yaml:"image"`
	Resources    ResourceConfig  `yaml:"resources"`
	Ports        []ContainerPort `yaml:"ports"`
	Entrance     string          `yaml:"entrance"`
	VolumeMounts []VolumeMount   `yaml:"volumeMounts,omitempty"`
	LogDirectory string          `yaml:"logDirectory,omitempty"`
}

// ContainerPort 表示容器端口配置
type ContainerPort struct {
	ContainerPort int `yaml:"containerPort"`
}

// ResourceConfig 表示资源限制配置
type ResourceConfig struct {
	Requests ResourceRequest `yaml:"requests"`
	Limits   ResourceLimit   `yaml:"limits"`
}

// ResourceRequest 表示资源请求
type ResourceRequest struct {
	Memory string `yaml:"memory"`
}

// ResourceLimit 表示资源限制
type ResourceLimit struct {
	Memory string `yaml:"memory"`
}

// VolumeMount 表示卷挂载配置
type VolumeMount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

// GetContainerNames 获取所有容器名称
func (c *ComposeConfig) GetContainerNames() []string {
	var names []string
	for _, container := range c.Service.Containers {
		names = append(names, container.Name)
	}
	return names
}

// GetContainerImages 获取所有容器镜像
func (c *ComposeConfig) GetContainerImages() []string {
	var images []string
	for _, container := range c.Service.Containers {
		images = append(images, container.Image)
	}
	return images
}

// GetContainerPorts 获取指定容器的端口
func (c *ComposeConfig) GetContainerPorts(containerName string) []int {
	for _, container := range c.Service.Containers {
		if container.Name == containerName {
			var ports []int
			for _, port := range container.Ports {
				ports = append(ports, port.ContainerPort)
			}
			return ports
		}
	}
	return nil
}

// LoadAndValidateComposeFile 加载并验证compose.yaml文件
func LoadAndValidateComposeFile(filePath string) (*ComposeConfig, error) {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read compose file: %v", err)
	}

	// 解析YAML
	config, err := ParseComposeYAML(content)
	if err != nil {
		return nil, err
	}

	// 验证配置
	if err := ValidateComposeConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// ParseComposeYAML 解析compose.yaml文件内容
func ParseComposeYAML(content []byte) (*ComposeConfig, error) {
	var config ComposeConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}
	return &config, nil
}

// ValidateComposeConfig 验证compose配置
func ValidateComposeConfig(config *ComposeConfig) error {
	if config == nil {
		return fmt.Errorf("compose config is nil")
	}

	// 验证service
	if err := validateServiceConfig(&config.Service); err != nil {
		return fmt.Errorf("service validation failed: %v", err)
	}

	return nil
}

// validateServiceConfig 验证service配置
func validateServiceConfig(service *ServiceConfig) error {
	if service == nil {
		return fmt.Errorf("service config is nil")
	}

	// 验证containers不能为空
	if len(service.Containers) == 0 {
		return fmt.Errorf("containers cannot be empty")
	}

	// 验证每个容器
	for i, container := range service.Containers {
		if err := validateContainer(&container, i); err != nil {
			return err
		}
	}

	return nil
}

// validateContainer 验证容器配置
func validateContainer(container *Container, index int) error {
	if container == nil {
		return fmt.Errorf("container at index %d is nil", index)
	}

	// 验证name不能为空
	if strings.TrimSpace(container.Name) == "" {
		return fmt.Errorf("container[%d].name cannot be empty", index)
	}

	// 验证image不能为空
	if strings.TrimSpace(container.Image) == "" {
		return fmt.Errorf("container[%d].image cannot be empty", index)
	}

	if !isValidImageName(container.Image) {
		return fmt.Errorf("container[%s].image has invalid format: %s", container.Name, container.Image)
	}

	if len(container.Ports) > 0 {
		// 验证每个端口
		for j, port := range container.Ports {
			if port.ContainerPort <= 0 {
				return fmt.Errorf("container[%d].ports[%d].containerPort must be greater than 0", index, j)
			}
		}
	}

	// 验证resources
	if err := validateResourceConfig(&container.Resources, index); err != nil {
		return err
	}

	return nil
}

// validateResourceConfig 验证资源配置
func validateResourceConfig(resources *ResourceConfig, containerIndex int) error {
	if resources == nil {
		return fmt.Errorf("container[%d].resources cannot be empty", containerIndex)
	}

	// 验证requests
	if strings.TrimSpace(resources.Requests.Memory) == "" {
		return fmt.Errorf("container[%d].resources.requests.memory cannot be empty", containerIndex)
	}

	// 验证limits
	if strings.TrimSpace(resources.Limits.Memory) == "" {
		return fmt.Errorf("container[%d].resources.limits.memory cannot be empty", containerIndex)
	}

	return nil
}

// isValidImageName 验证镜像名称格式
func isValidImageName(image string) bool {
	// 基本格式检查：应该包含冒号（tag）或斜杠（registry/repository）
	if strings.Contains(image, ":") {
		// 有tag的镜像
		parts := strings.Split(image, ":")
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			return false
		}
	} else if strings.Contains(image, "/") {
		// 没有tag但有registry/repository
		return true
	}
	// 简单的镜像名（如nginx）也是有效的
	return strings.TrimSpace(image) != ""
}
