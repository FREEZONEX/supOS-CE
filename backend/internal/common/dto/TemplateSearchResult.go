package dto

// TemplateSearchResult 模板搜索结果
type TemplateSearchResult struct {
	// ID 模板ID
	ID string `gorm:"column:id" json:"id" schema:"Id"`

	// Name 模板名称
	Name string `gorm:"column:name" json:"name" schema:"模板名称"`

	// Description 模型描述
	Description string `gorm:"column:description" json:"description" schema:"模型描述"`

	// Alias 别名
	Alias string `gorm:"column:alias" json:"alias" schema:"别名"`
}
