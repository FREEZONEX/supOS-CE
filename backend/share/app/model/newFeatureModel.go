package model

import (
	"backend/internal/common/errors"
	"backend/share/app/util"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// NewFeatureModel 新功能模型
// 用于描述和管理应用的新功能特性
type NewFeatureModel struct {
	// Name 功能名称
	Name string `json:"name" yaml:"name"`

	// Description 功能描述
	Description string `json:"description" yaml:"description"`

	// ImagePath 镜像文件路径（本地文件路径）
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
		return fmt.Errorf("app name can't be empty")
	}

	// 验证非法字符
	if !ValidateFilenameBasic(m.Name) {
		return fmt.Errorf("app name contains invalid characters")
	}
	// 验证该名称是否已经存在
	fullPath := filepath.Join(util.AppInstalledDir, m.Name)
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("app name has been used")
	}

	// 镜像路径和URL至少需要提供一个
	if m.ImagePath == "" && m.ImageUrl == "" {
		return fmt.Errorf("镜像路径和镜像URL至少需要提供一个")
	}

	// 如果提供了镜像路径，验证文件是否存在
	if m.ImagePath != "" {
		if _, err := os.Stat(m.ImagePath); os.IsNotExist(err) {
			return fmt.Errorf("镜像文件不存在: %s", m.ImagePath)
		}

		// 验证文件格式
		ext := filepath.Ext(m.ImagePath)
		if ext != ".tar" && ext != ".tar.gz" && ext != ".tgz" {
			return fmt.Errorf("不支持的镜像文件格式: %s，仅支持 .tar, .tar.gz, .tgz 格式", ext)
		}
	}

	// 验证菜单URL格式（如果提供了菜单URL）
	if m.Menu == nil {
		return errors.NewAppErrorWithMsg("menu can't be empty")
	}
	if err := m.Menu.Validate(); err != nil {
		return err
	}

	// 验证Docker Compose配置（如果提供了）
	if m.ComposeYaml != "" {
		// 简单的YAML格式验证
		if !strings.Contains(m.ComposeYaml, "version:") && !strings.Contains(m.ComposeYaml, "services:") {
			return fmt.Errorf("Docker Compose配置格式不正确，应包含 'version' 和 'services' 字段")
		}
	}

	return nil
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
		return nil, fmt.Errorf("配置文件不存在: %s", filePath)
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	var model NewFeatureModel
	if err := yaml.Unmarshal(data, &model); err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %v", err)
	}

	// 验证模型
	if err := model.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	return &model, nil
}

// SaveToFile 将模型保存到YAML文件
func (m *NewFeatureModel) SaveToFile(filePath string) error {
	// 验证模型
	if err := m.Validate(); err != nil {
		return fmt.Errorf("模型验证失败: %v", err)
	}

	// 转换为YAML
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("序列化YAML失败: %v", err)
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
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
