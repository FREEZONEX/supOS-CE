package model

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// AppConfig 表示app.yaml文件的完整结构
type AppConfig struct {
	APIVersion  string      `yaml:"apiVersion"`
	Name        string      `yaml:"name"`
	VendorName  string      `yaml:"vendorName"`
	AppVersion  string      `yaml:"appVersion"`
	ShowName    string      `yaml:"showName"`
	Description string      `yaml:"description"`
	Doc         interface{} `yaml:"doc,omitempty"`
	Icon        string      `yaml:"icon,omitempty"`
	IndexUrl    string      `yaml:"indexUrl"`
	Type        string      `yaml:"type"`
	DeployMode  string      `yaml:"deployMode"`
	AppID       string      `yaml:"appId,omitempty"`
}

// Validate 验证app配置
func (c *AppConfig) Validate() error {
	// 验证name不能为空
	if strings.TrimSpace(c.Name) == "" {
		return NewValidationError("name cannot be empty")
	}

	// 验证vendorName不能为空
	if strings.TrimSpace(c.VendorName) == "" {
		return NewValidationError("vendorName cannot be empty")
	}

	// 验证appVersion不能为空
	if strings.TrimSpace(c.AppVersion) == "" {
		return NewValidationError("appVersion cannot be empty")
	}

	// 验证showName不能为空
	if strings.TrimSpace(c.ShowName) == "" {
		return NewValidationError("showName cannot be empty")
	}

	// 验证indexUrl不能为空
	if strings.TrimSpace(c.IndexUrl) == "" {
		return NewValidationError("indexUrl cannot be empty")
	}

	return nil
}

// GetAppInfo 获取应用基本信息
func (c *AppConfig) GetAppInfo() map[string]string {
	return map[string]string{
		"name":       c.Name,
		"vendorName": c.VendorName,
		"appVersion": c.AppVersion,
		"showName":   c.ShowName,
		"type":       c.Type,
		"deployMode": c.DeployMode,
		"indexUrl":   c.IndexUrl,
	}
}

// LoadAndValidateAppConfig 从文件加载并验证app配置
// filePath: app.yaml文件路径
// 返回: 解析后的AppConfig和错误信息
func LoadAndValidateAppConfig(filePath string) (*AppConfig, error) {
	// 1. 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, &ValidationError{Message: fmt.Sprintf("failed to read app.yaml file: %v", err)}
	}

	// 2. 解析YAML
	var config AppConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, &ValidationError{Message: fmt.Sprintf("failed to parse app.yaml: %v", err)}
	}

	// 3. 验证配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}
