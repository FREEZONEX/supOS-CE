package relationDB

import "time"

// suposResource maps to table supos_resource.
type SuposResource struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ParentID        *int64    `gorm:"column:parent_id"`
	Type            int       `gorm:"column:type"`
	Source          *string   `gorm:"column:source"`
	Code            string    `gorm:"column:code"`
	NameCode        *string   `gorm:"column:name_code"`
	RouteSource     *int      `gorm:"column:route_source"`
	URL             *string   `gorm:"column:url"`
	URLType         *int      `gorm:"column:url_type"`
	OpenType        *int      `gorm:"column:open_type"`
	Icon            *string   `gorm:"column:icon"`
	DescriptionCode *string   `gorm:"column:description_code"`
	Sort            *int      `gorm:"column:sort"`
	EditEnable      *bool     `gorm:"column:edit_enable"`
	HomeEnable      *bool     `gorm:"column:home_enable"`
	Fixed           *bool     `gorm:"column:fixed"`
	Enable          *bool     `gorm:"column:enable"`
	UpdateAt        time.Time `gorm:"column:update_at"`
	CreateAt        time.Time `gorm:"column:create_at"`
}

func (SuposResource) TableName() string {
	return "supos_resource"
}

// supos_i18n_language 国际化语言表
type SuposI18nLanguage struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"` // 主键
	LanguageCode string `gorm:"column:language_code"`               // 语言编码，例如 zh_CN、en_US
	LanguageName string `gorm:"column:language_name"`               // 语言名称，例如 中文（中国）、English (US)
	LanguageType int    `gorm:"column:language_type"`               // 语言类型：1-内置，2-自定义
	HasUsed      bool   `gorm:"column:has_used"`                    // 是否启用
}

func (SuposI18nLanguage) TableName() string {
	return "supos_i18n_language"
}
