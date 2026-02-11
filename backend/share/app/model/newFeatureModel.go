package model

import (
	"backend/share/app/util"
	"fmt"
	"os"
	"path/filepath"
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
		// 简单的YAML格式验证
		if !strings.Contains(m.ComposeYaml, "services:") {
			return errors.Parameter.WithMsg(fmt.Sprintf("app.compose.content.empty"))
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

	// 验证模型
	if err := model.Validate(); err != nil {
		return nil, errors.Parameter.WithMsg(fmt.Sprintf("配置验证失败: %v", err))
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
