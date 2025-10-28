package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// FieldType represents the data type of a field.
type FieldType string

const (
	FieldTypeInteger  = "INTEGER"
	FieldTypeLong     = "LONG"
	FieldTypeFloat    = "FLOAT"
	FieldTypeDouble   = "DOUBLE"
	FieldTypeBoolean  = "BOOLEAN"
	FieldTypeDatetime = "DATETIME"
	FieldTypeString   = "STRING"
	FieldTypeBlob     = "BLOB"
	FieldTypeLBlob    = "LBLOB"
)

// fieldTypeInfo holds the metadata for each FieldType constant.
type fieldTypeInfo struct {
	isNumber     bool
	defaultValue any
}

// fieldTypeDetails maps each FieldType constant to its metadata.
var fieldTypeDetails = map[FieldType]fieldTypeInfo{
	FieldTypeInteger:  {true, 0},
	FieldTypeLong:     {true, int64(0)},
	FieldTypeFloat:    {true, float32(0.0)},
	FieldTypeDouble:   {true, float64(0.0)},
	FieldTypeBoolean:  {false, false},
	FieldTypeDatetime: {false, nil},
	FieldTypeString:   {false, nil},
	FieldTypeBlob:     {false, nil},
	FieldTypeLBlob:    {false, nil},
}

var fieldTypes []FieldType

// init populates the nameMap for fast lookups.
func init() {
	fieldTypes = make([]FieldType, 0, len(fieldTypeDetails))
	for ft := range fieldTypeDetails {
		fieldTypes = append(fieldTypes, ft)
	}
	sort.Slice(fieldTypes, func(i, j int) bool {
		return fieldTypes[i] < fieldTypes[j]
	})
}

// Name returns the canonical string name of the field type.
func (f FieldType) Name() string {
	return string(f)
}

// IsNumber returns true if the field type is numeric.
func (f FieldType) IsNumber() bool {
	switch f {
	case FieldTypeInteger, FieldTypeLong, FieldTypeFloat, FieldTypeDouble:
		return true
	}
	return false
}

// DefaultValue returns the default value for the field type.
func (f FieldType) DefaultValue() any {
	return fieldTypeDetails[f].defaultValue
}

// String implements the fmt.Stringer interface for easy printing.
func (f FieldType) String() string {
	return f.Name()
}

func FieldTypes() (ts []FieldType) {
	return fieldTypes
}

func GetFieldTypeByName(name string) (FieldType, bool) {
	return GetFieldTypeByNameIgnoreCase(name)
}

// GetFieldTypeByNameIgnoreCase finds a FieldType by its name, case-insensitively.
// It includes a special case to handle "int" as an alias for Integer.
func GetFieldTypeByNameIgnoreCase(name string) (FieldType, bool) {
	ft := FieldType(strings.ToUpper(name))
	_, ok := fieldTypeDetails[ft]
	if ok {
		return ft, true
	}

	if strings.EqualFold(name, "int") {
		return FieldTypeInteger, true
	}

	return "", false
}

// MarshalJSON implements the json.Marshaler interface, serializing the FieldType to its string name.
func (f FieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Name())
}

// UnmarshalJSON implements the json.Unmarshaler interface, deserializing a string into a FieldType.
func (f *FieldType) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	ft, ok := GetFieldTypeByNameIgnoreCase(name)
	if !ok {
		return fmt.Errorf("invalid FieldType name: %s", name)
	}

	*f = ft
	return nil
}
