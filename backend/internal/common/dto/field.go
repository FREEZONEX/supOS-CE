package dto

import (
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"strings"
)

// FieldDefine represents field definition
type FieldDefine struct {
	Name        string          `json:"name" validate:"required"` // 字段名
	Type        enums.FieldType `json:"type" validate:"required"` // 字段类型
	Unique      bool            `json:"unique,omitzero"`          // 是否唯一约束
	Index       string          `json:"index,omitzero"`           // 对应的协议字段key
	DisplayName string          `json:"displayName,omitzero"`     // 显式名
	Remark      string          `json:"remark,omitzero"`          // 备注
	MaxLen      int             `json:"maxLen,omitzero"`          // 最大长度(string字段类型生效)
	TbValueName string          `json:"-"`                        // Internal use
	Unit        string          `json:"unit,omitzero"`            // 位号单位
	UpperLimit  float64         `json:"upperLimit,omitzero"`      // 原始上限
	LowerLimit  float64         `json:"lowerLimit,omitzero"`      // 原始下限
	Decimal     int             `json:"decimal,omitzero"`         // 小数精度位数
}

// IsUnique checks if field has unique constraint
func (f *FieldDefine) IsUnique() bool {
	return f.Unique
}

// IsSystemField checks if field is a system field
func (f *FieldDefine) IsSystemField() bool {
	_, ok := constants.SystemFields[f.Name]
	return strings.HasPrefix(f.Name, constants.SystemFieldPrev) || ok
}

// SetName sets and trims the field name
func (f *FieldDefine) SetName(name string) {
	f.Name = strings.TrimSpace(name)
}

// SetIndex sets and trims the field index
func (f *FieldDefine) SetIndex(index string) {
	f.Index = strings.TrimSpace(index)
}

// Clone creates a deep copy of FieldDefine
func (f *FieldDefine) Clone() *FieldDefine {
	clone := &FieldDefine{
		Name:        f.Name,
		Type:        f.Type,
		Index:       f.Index,
		DisplayName: f.DisplayName,
		Remark:      f.Remark,
		TbValueName: f.TbValueName,
		Unit:        f.Unit,
		UpperLimit:  f.UpperLimit,
		LowerLimit:  f.LowerLimit,
		Decimal:     f.Decimal,
		Unique:      f.Unique,
		MaxLen:      f.MaxLen,
	}

	return clone
}

// FieldDefines represents a collection of field definitions
type FieldDefines struct {
	FieldsMap     map[string]*FieldDefine // Field name -> FieldDefine
	FieldIndexMap map[string]string       // Index -> Field name
	UniqueKeys    map[string]bool         // Set of unique field names
	CalcField     *FieldDefine            // Calculation field
}

// NewFieldDefines creates a new FieldDefines from a slice
func NewFieldDefines(fields []*FieldDefine) *FieldDefines {
	fd := &FieldDefines{
		FieldsMap:     make(map[string]*FieldDefine),
		FieldIndexMap: make(map[string]string),
		UniqueKeys:    make(map[string]bool),
	}

	if len(fields) > 0 {
		for _, f := range fields {
			fd.FieldsMap[f.Name] = f

			if f.IsUnique() {
				fd.UniqueKeys[f.Name] = true
			}

			if f.Index != "" {
				fd.FieldIndexMap[f.Index] = f.Name
			}
		}
	}

	return fd
}

// NewFieldDefinesFromMap creates a new FieldDefines from a map
func NewFieldDefinesFromMap(fieldsMap map[string]*FieldDefine) *FieldDefines {
	if len(fieldsMap) == 0 {
		return &FieldDefines{
			FieldsMap:     make(map[string]*FieldDefine),
			FieldIndexMap: make(map[string]string),
			UniqueKeys:    make(map[string]bool),
		}
	}

	fields := make([]*FieldDefine, 0, len(fieldsMap))
	for _, f := range fieldsMap {
		fields = append(fields, f)
	}

	return NewFieldDefines(fields)
}

// ToFieldDefineArray converts to array
func (fd *FieldDefines) ToFieldDefineArray() []*FieldDefine {
	if len(fd.FieldsMap) == 0 {
		return []*FieldDefine{}
	}

	fields := make([]*FieldDefine, 0, len(fd.FieldsMap))
	for _, f := range fd.FieldsMap {
		fields = append(fields, f)
	}

	return fields
}

// DefaultMaxStrLen is the default maximum string length
const DefaultMaxStrLen = 512

// InstanceField represents a reference to a field in a UNS instance
type InstanceField struct {
	ID    int64  `json:"id,omitzero"`    // UNS ID
	Alias string `json:"alias,omitzero"` // Alias (filled by backend)
	Path  string `json:"path,omitzero"`  // Path (filled by backend)
	Field string `json:"field,omitzero"` // Field name
	UTS   bool   `json:"uts,omitzero"`   // true--计算型实例，使用当前uns的时间戳
}

// GetTopic returns the topic based on configuration
func (i *InstanceField) GetTopic() string {
	if constants.UseAliasAsTopic {
		return i.Alias
	}
	return i.Path
}

// NewInstanceField creates a new InstanceField with ID and field
func NewInstanceField(id int64, field string) *InstanceField {
	return &InstanceField{
		ID:    id,
		Field: field,
	}
}

// NewInstanceFieldWithAlias creates a new InstanceField with alias and field
func NewInstanceFieldWithAlias(alias, field string) *InstanceField {
	return &InstanceField{
		Alias: alias,
		Field: field,
	}
}

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
