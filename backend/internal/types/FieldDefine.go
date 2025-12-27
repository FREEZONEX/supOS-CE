package types

import (
	"backend/internal/common/constants"
	"strings"
)

type FieldDefineInfo interface {
	GetName() string
	GetType() string
	GetUnique() *bool
	GetIndex() *string
	GetDisplayName() *string
	GetRemark() *string
	GetMaxLen() *int
	GetTbValueName() *string
	GetUnit() *string
	GetUpperLimit() *float64
	GetLowerLimit() *float64
	GetDecimal() *int
}

const DefaultMaxStrLen = 512

// IsUnique checks if field has unique constraint
func (f *FieldDefine) IsUnique() bool {
	return f.Unique != nil && *f.Unique
}

// IsSystemField checks if field is a system field
func (f *FieldDefine) IsSystemField() bool {
	_, ok := constants.SystemFields[f.Name]
	sf := f.SystemField
	return ok || (sf != nil && *sf) || strings.HasPrefix(f.Name, constants.SystemFieldPrev)
}

// SetName sets and trims the field name
func (f *FieldDefine) SetName(name string) {
	f.Name = strings.TrimSpace(name)
}

// SetIndex sets and trims the field index
func (f *FieldDefine) SetIndex(index string) {
	index = strings.TrimSpace(index)
	if len(index) > 0 {
		f.Index = &index
	}
}
func (f *FieldDefine) Equals(a *FieldDefine) bool {
	return f.Name == a.Name && f.Type == a.Type
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

// GetName 获取 Name
func (f *FieldDefine) GetName() string {
	return f.Name
}

// GetType 获取 Type
func (f *FieldDefine) GetType() FieldType {
	return FieldType(f.Type)
}

// GetUnique 获取 Unique
func (f *FieldDefine) GetUnique() *bool {
	return f.Unique
}

// GetIndex 获取 Index
func (f *FieldDefine) GetIndex() *string {
	return f.Index
}

// GetDisplayName 获取 DisplayName
func (f *FieldDefine) GetDisplayName() *string {
	return f.DisplayName
}

// GetRemark 获取 Remark
func (f *FieldDefine) GetRemark() *string {
	return f.Remark
}

// GetMaxLen 获取 MaxLen
func (f *FieldDefine) GetMaxLen() (maxLen *int) {
	return f.MaxLen
}

// GetTbValueName 获取 TbValueName
func (f *FieldDefine) GetTbValueName() *string {
	return f.TbValueName
}

// GetUnit 获取 Unit
func (f *FieldDefine) GetUnit() *string {
	return f.Unit
}

// GetUpperLimit 获取 UpperLimit
func (f *FieldDefine) GetUpperLimit() *float64 {
	return f.UpperLimit
}

// GetLowerLimit 获取 LowerLimit
func (f *FieldDefine) GetLowerLimit() *float64 {
	return f.LowerLimit
}

// GetDecimal 获取 Decimal
func (f *FieldDefine) GetDecimal() *int {
	return f.Decimal
}

var UnsLastValueFill func(uns *CreateTopicDto)

func (f *FieldDefine) GetLastValue() interface{} {
	if f.LastTime == 0 {
		f.tryFillLastValue()
	}
	return f.LastValue
}
func (f *FieldDefine) GetLastTime() int64 {
	if f.LastTime == 0 {
		f.tryFillLastValue()
	}
	return f.LastTime
}
func (f *FieldDefine) tryFillLastValue() {
	if uns, ok := f.Uns.(*CreateTopicDto); ok && UnsLastValueFill != nil {
		UnsLastValueFill(uns)
		if f.LastTime == 0 {
			f.LastTime = -1
		}
	}
}
