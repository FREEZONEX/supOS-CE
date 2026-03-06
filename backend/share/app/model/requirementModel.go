package model

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// RequirementConfig 表示requirement.yaml文件的完整结构
type RequirementConfig struct {
	APIVersion   string       `yaml:"apiVersion"`
	Requirements Requirements `yaml:"requirements"`
}

// Requirements 表示requirements配置
type Requirements struct {
	Databases []DatabaseRequirement `yaml:"databases,omitempty"`
	DPS       []DPSRequirement      `yaml:"dps,omitempty"`
	Volumes   []VolumeRequirement   `yaml:"volumes,omitempty"`
}

// DatabaseRequirement 表示数据库需求
type DatabaseRequirement struct {
	Name         string `yaml:"name"`
	ResourceType string `yaml:"resourceType"`
	Apply        string `yaml:"apply"`
}

// DPSRequirement 表示DPS需求
type DPSRequirement struct {
	Name string `yaml:"name"`
}

// VolumeRequirement 表示卷需求
type VolumeRequirement struct {
	Name         string `yaml:"name"`
	ResourceType string `yaml:"resourceType"`
	Size         string `yaml:"size"`
	LocalPath    string `yaml:"localPath"`
}

// Validate 验证requirement配置
func (c *RequirementConfig) Validate() error {
	if c == nil {
		return NewValidationError("requirement config is nil")
	}

	// 验证apiVersion不能为空
	if strings.TrimSpace(c.APIVersion) == "" {
		return NewValidationError("apiVersion cannot be empty")
	}

	// 验证requirements
	if err := c.Requirements.Validate(); err != nil {
		return fmt.Errorf("requirements validation failed: %v", err)
	}

	return nil
}

// Validate 验证requirements配置
func (r *Requirements) Validate() error {
	if r == nil {
		return NewValidationError("requirements is nil")
	}

	// 验证databases
	for i, db := range r.Databases {
		if err := db.Validate(i); err != nil {
			return err
		}
	}

	// 验证dps
	for i, dps := range r.DPS {
		if err := dps.Validate(i); err != nil {
			return err
		}
	}

	// 验证volumes
	for i, volume := range r.Volumes {
		if err := volume.Validate(i); err != nil {
			return err
		}
	}

	return nil
}

// Validate 验证数据库需求
func (d *DatabaseRequirement) Validate(index int) error {
	// 验证name不能为空
	if strings.TrimSpace(d.Name) == "" {
		return NewValidationError(fmt.Sprintf("databases[%d].name cannot be empty", index))
	}

	// 验证resourceType不能为空
	if strings.TrimSpace(d.ResourceType) == "" {
		return NewValidationError(fmt.Sprintf("databases[%d].resourceType cannot be empty", index))
	}

	// 验证apply不能为空
	if strings.TrimSpace(d.Apply) == "" {
		return NewValidationError(fmt.Sprintf("databases[%d].apply cannot be empty", index))
	}

	return nil
}

// Validate 验证DPS需求
func (d *DPSRequirement) Validate(index int) error {
	// 验证name不能为空
	if strings.TrimSpace(d.Name) == "" {
		return NewValidationError(fmt.Sprintf("dps[%d].name cannot be empty", index))
	}

	return nil
}

// Validate 验证卷需求
func (v *VolumeRequirement) Validate(index int) error {
	// 验证name不能为空
	if strings.TrimSpace(v.Name) == "" {
		return NewValidationError(fmt.Sprintf("volumes[%d].name cannot be empty", index))
	}

	// 验证resourceType不能为空
	if strings.TrimSpace(v.ResourceType) == "" {
		return NewValidationError(fmt.Sprintf("volumes[%d].resourceType cannot be empty", index))
	}

	// 验证size不能为空
	if strings.TrimSpace(v.Size) == "" {
		return NewValidationError(fmt.Sprintf("volumes[%d].size cannot be empty", index))
	}

	// 验证localPath不能为空
	if strings.TrimSpace(v.LocalPath) == "" {
		return NewValidationError(fmt.Sprintf("volumes[%d].localPath cannot be empty", index))
	}

	return nil
}

// GetDatabaseNames 获取所有数据库名称
func (r *Requirements) GetDatabaseNames() []string {
	var names []string
	for _, db := range r.Databases {
		names = append(names, db.Name)
	}
	return names
}

// GetDPSNames 获取所有DPS名称
func (r *Requirements) GetDPSNames() []string {
	var names []string
	for _, dps := range r.DPS {
		names = append(names, dps.Name)
	}
	return names
}

// GetVolumeNames 获取所有卷名称
func (r *Requirements) GetVolumeNames() []string {
	var names []string
	for _, volume := range r.Volumes {
		names = append(names, volume.Name)
	}
	return names
}

// GetVolumeByPath 根据路径获取卷
func (r *Requirements) GetVolumeByPath(path string) *VolumeRequirement {
	for _, volume := range r.Volumes {
		if volume.LocalPath == path {
			return &volume
		}
	}
	return nil
}

// HasDatabases 检查是否有数据库需求
func (r *Requirements) HasDatabases() bool {
	return len(r.Databases) > 0
}

// HasDPS 检查是否有DPS需求
func (r *Requirements) HasDPS() bool {
	return len(r.DPS) > 0
}

// HasVolumes 检查是否有卷需求
func (r *Requirements) HasVolumes() bool {
	return len(r.Volumes) > 0
}

// GetVolumeInfo 获取卷信息
func (v *VolumeRequirement) GetVolumeInfo() map[string]string {
	return map[string]string{
		"name":         v.Name,
		"resourceType": v.ResourceType,
		"size":         v.Size,
		"localPath":    v.LocalPath,
	}
}

// GetDatabaseInfo 获取数据库信息
func (d *DatabaseRequirement) GetDatabaseInfo() map[string]string {
	return map[string]string{
		"name":         d.Name,
		"resourceType": d.ResourceType,
		"apply":        d.Apply,
	}
}

// GetDPSInfo 获取DPS信息
func (d *DPSRequirement) GetDPSInfo() map[string]string {
	return map[string]string{
		"name": d.Name,
	}
}

// LoadAndValidateRequirementConfig 从文件加载并验证requirement配置
func LoadAndValidateRequirementConfig(filePath string) (*RequirementConfig, error) {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read requirement file: %v", err)
	}

	// 解析YAML
	var config RequirementConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse requirement YAML: %v", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}
