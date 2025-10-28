package types

// UpdateFieldDto represents field update DTO
type UpdateFieldDto struct {
	Alias     string         `json:"alias,omitzero"`     // 别名即表名
	Topic     string         `json:"topic,omitzero"`     // 主题
	NewFields []*FieldDefine `json:"newFields,omitzero"` // 新增的字段定义
	DelFields []*FieldDefine `json:"delFields,omitzero"` // 删除的字段定义
}

// NewUpdateFieldDto creates a new UpdateFieldDto
func NewUpdateFieldDto(alias, topic string, newFields, delFields []*FieldDefine) *UpdateFieldDto {
	return &UpdateFieldDto{
		Alias:     alias,
		Topic:     topic,
		NewFields: newFields,
		DelFields: delFields,
	}
}
