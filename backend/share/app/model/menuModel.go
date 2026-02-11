package model

import (
	"backend/internal/common/errors"
	"fmt"
	"strings"
)

// MenuModel 菜单模型
// 用于描述和管理系统菜单项
type MenuModel struct {

	// Name 菜单名称
	Name string `json:"name" yaml:"name" gorm:"column:name;not null;uniqueIndex"`

	// Description 菜单描述
	Description string `json:"description" yaml:"description" gorm:"column:description"`

	// IndexUrl 菜单索引URL
	// 菜单项对应的前端路由或页面URL
	IndexUrl string `json:"indexUrl" yaml:"indexUrl" gorm:"column:index_url"`

	// OpenType 打开类型
	// 0=当前窗口打开，1=新窗口打开，2=iframe嵌入，3=外部链接
	OpenType int `json:"openType" yaml:"openType" gorm:"column:open_type;default:0"`

	// StripPath 是否剥离路径
	// 当为true时，代理会从请求路径中剥离匹配的前缀
	StripPath bool `json:"stripPath" yaml:"stripPath" gorm:"column:strip_path;default:true"`

	// PreserveHost 是否保留主机头
	// 当为true时，代理会保留原始请求的主机头
	PreserveHost bool `json:"preserveHost" yaml:"preserveHost" gorm:"column:preserve_host;default:false"`

	// IconUrl 图标URL
	// 菜单图标地址，可以是本地路径或远程URL
	IconUrl string `json:"iconUrl" yaml:"iconUrl" gorm:"column:icon_url"`
}

// Validate 验证菜单模型字段
func (m *MenuModel) Validate() error {

	// 验证 Name 字段
	if err := validateName(m.Name); err != nil {
		return err
	}

	// 验证 IndexUrl 字段
	if err := validateIndexUrl(m.IndexUrl); err != nil {
		return err
	}

	// 验证 IconUrl 字段（如果提供）
	if m.IconUrl != "" {
		if err := validateIconUrl(m.IconUrl); err != nil {
			return err
		}
	}

	// 验证 Description 字段长度
	if err := validateDescription(m.Description); err != nil {
		return err
	}

	return nil
}

// ValidateForCreate 验证创建菜单时的字段
func (m *MenuModel) ValidateForCreate() error {
	// 首先执行基本验证
	if err := m.Validate(); err != nil {
		return err
	}

	// 创建时 Name 必须非空
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("菜单名称不能为空")
	}

	return nil
}

// ValidateForUpdate 验证更新菜单时的字段
func (m *MenuModel) ValidateForUpdate() error {
	// 更新时执行基本验证
	return m.Validate()
}

// IsValidUrl 检查URL是否有效
func (m *MenuModel) IsValidUrl() bool {
	if m.IndexUrl == "" {
		return false
	}

	// 检查是否是有效的URL或路径
	if strings.HasPrefix(m.IndexUrl, "http://") ||
		strings.HasPrefix(m.IndexUrl, "https://") ||
		strings.HasPrefix(m.IndexUrl, "/") {
		return true
	}

	return false
}

// IsExternalLink 检查是否是外部链接
func (m *MenuModel) IsExternalLink() bool {
	if m.IndexUrl == "" {
		return false
	}

	return strings.HasPrefix(m.IndexUrl, "http://") ||
		strings.HasPrefix(m.IndexUrl, "https://")
}

// IsInternalPath 检查是否是内部路径
func (m *MenuModel) IsInternalPath() bool {
	if m.IndexUrl == "" {
		return false
	}

	return strings.HasPrefix(m.IndexUrl, "/")
}

// Helper validation functions

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.NewAppErrorWithMsg("menu name can't be empty")
	}

	// 名称长度限制
	if len(name) > 100 {
		return fmt.Errorf("菜单名称长度不能超过100个字符")
	}

	// 名称格式检查（只允许字母、数字、下划线、中划线）
	for _, r := range name {
		if !isValidNameChar(r) {
			return fmt.Errorf("菜单名称只能包含字母、数字、下划线和中划线")
		}
	}

	return nil
}

func validateDisplayName(displayName string) error {
	if strings.TrimSpace(displayName) == "" {
		return fmt.Errorf("菜单显示名称不能为空")
	}

	// 显示名称长度限制
	if len(displayName) > 200 {
		return fmt.Errorf("菜单显示名称长度不能超过200个字符")
	}

	return nil
}

func validateIndexUrl(indexUrl string) error {
	if strings.TrimSpace(indexUrl) == "" {
		return fmt.Errorf("菜单索引URL不能为空")
	}

	// 检查URL格式
	if !isValidUrlOrPath(indexUrl) {
		return fmt.Errorf("无效的菜单索引URL格式，必须是有效的URL或路径")
	}

	// URL长度限制
	if len(indexUrl) > 500 {
		return fmt.Errorf("菜单索引URL长度不能超过500个字符")
	}

	return nil
}

func validateIconUrl(iconUrl string) error {
	// 图标URL可以为空
	if iconUrl == "" {
		return nil
	}

	// 检查图标URL格式
	if !isValidUrlOrPath(iconUrl) {
		return fmt.Errorf("无效的图标URL格式，必须是有效的URL或路径")
	}

	// URL长度限制
	if len(iconUrl) > 500 {
		return fmt.Errorf("图标URL长度不能超过500个字符")
	}

	// 检查文件扩展名（如果是本地文件）
	if strings.HasPrefix(iconUrl, "/") {
		ext := strings.ToLower(getFileExtension(iconUrl))
		if ext != "" {
			allowedExtensions := []string{".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico"}
			found := false
			for _, allowed := range allowedExtensions {
				if ext == allowed {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("不支持的图标文件格式: %s，支持的格式: %v", ext, allowedExtensions)
			}
		}
	}

	return nil
}

func validateDescription(description string) error {
	// 描述可以为空
	if description == "" {
		return nil
	}

	// 描述长度限制
	if len(description) > 200 {
		return fmt.Errorf("菜单描述长度不能超过200个字符")
	}

	return nil
}

// Utility functions

func isValidNameChar(r rune) bool {
	// 允许：字母、数字、下划线、中划线、中文
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_' || r == '-' ||
		(r >= 0x4E00 && r <= 0x9FFF) // 中文范围
}

func isValidUrlOrPath(urlStr string) bool {
	// 检查是否是有效的URL或路径
	if strings.HasPrefix(urlStr, "http://") ||
		strings.HasPrefix(urlStr, "https://") ||
		strings.HasPrefix(urlStr, "/") {
		return true
	}

	// 对于相对路径，也需要检查
	if strings.Contains(urlStr, ".") && !strings.Contains(urlStr, " ") {
		return true
	}

	return false
}

func getFileExtension(filename string) string {
	if filename == "" {
		return ""
	}

	// 查找最后一个点
	lastDot := strings.LastIndex(filename, ".")
	if lastDot == -1 || lastDot == len(filename)-1 {
		return ""
	}

	return filename[lastDot:]
}

// NewMenuModel 创建新的菜单模型（带验证）
func NewMenuModel(name, displayName, indexUrl string, openType int) (*MenuModel, error) {
	menu := &MenuModel{
		Name:         name,
		IndexUrl:     indexUrl,
		OpenType:     openType,
		StripPath:    true,  // 默认剥离路径
		PreserveHost: false, // 默认不保留主机头
	}

	if err := menu.ValidateForCreate(); err != nil {
		return nil, err
	}

	return menu, nil
}

// Clone 克隆菜单对象
func (m *MenuModel) Clone() *MenuModel {
	if m == nil {
		return nil
	}

	return &MenuModel{
		Name:         m.Name,
		Description:  m.Description,
		IndexUrl:     m.IndexUrl,
		OpenType:     m.OpenType,
		StripPath:    m.StripPath,
		PreserveHost: m.PreserveHost,
		IconUrl:      m.IconUrl,
	}
}

// ToMap 转换为map
func (m *MenuModel) ToMap() map[string]interface{} {
	if m == nil {
		return nil
	}

	return map[string]interface{}{
		"name":         m.Name,
		"description":  m.Description,
		"indexUrl":     m.IndexUrl,
		"openType":     m.OpenType,
		"stripPath":    m.StripPath,
		"preserveHost": m.PreserveHost,
		"iconUrl":      m.IconUrl,
	}
}

// FromMap 从map加载数据
func (m *MenuModel) FromMap(data map[string]interface{}) error {
	if m == nil {
		return fmt.Errorf("菜单模型不能为空")
	}

	if data == nil {
		return fmt.Errorf("数据不能为空")
	}

	// 从map中提取数据
	if val, ok := data["name"].(string); ok {
		m.Name = val
	}
	if val, ok := data["description"].(string); ok {
		m.Description = val
	}
	if val, ok := data["indexUrl"].(string); ok {
		m.IndexUrl = val
	}
	if val, ok := data["openType"].(int); ok {
		m.OpenType = val
	} else if val, ok := data["openType"].(float64); ok {
		m.OpenType = int(val)
	}
	if val, ok := data["stripPath"].(bool); ok {
		m.StripPath = val
	}
	if val, ok := data["preserveHost"].(bool); ok {
		m.PreserveHost = val
	}
	if val, ok := data["iconUrl"].(string); ok {
		m.IconUrl = val
	}

	// 验证数据
	return m.Validate()
}
