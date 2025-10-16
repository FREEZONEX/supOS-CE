package enums

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldType represents the data type of a field.
type FieldType int

const (
	FieldTypeInteger FieldType = iota
	FieldTypeLong
	FieldTypeFloat
	FieldTypeDouble
	FieldTypeBoolean
	FieldTypeDatetime
	FieldTypeString
	FieldTypeBlob
	FieldTypeLBlob
)

// fieldTypeInfo holds the metadata for each FieldType constant.
type fieldTypeInfo struct {
	name         string
	isNumber     bool
	defaultValue any
}

// fieldTypeDetails maps each FieldType constant to its metadata.
var fieldTypeDetails = map[FieldType]fieldTypeInfo{
	FieldTypeInteger:  {"INTEGER", true, 0},
	FieldTypeLong:     {"LONG", true, int64(0)},
	FieldTypeFloat:    {"FLOAT", true, float32(0.0)},
	FieldTypeDouble:   {"DOUBLE", true, float64(0.0)},
	FieldTypeBoolean:  {"BOOLEAN", false, false},
	FieldTypeDatetime: {"DATETIME", false, nil},
	FieldTypeString:   {"STRING", false, nil},
	FieldTypeBlob:     {"BLOB", false, nil},
	FieldTypeLBlob:    {"LBLOB", false, nil},
}

// nameMap provides a fast, case-insensitive lookup from a string name to a FieldType.
var nameMap = make(map[string]FieldType)

// init populates the nameMap for fast lookups.
func init() {
	for ft, info := range fieldTypeDetails {
		nameMap[strings.ToUpper(info.name)] = ft
	}
}

// Name returns the canonical string name of the field type.
func (f FieldType) Name() string {
	return fieldTypeDetails[f].name
}

// IsNumber returns true if the field type is numeric.
func (f FieldType) IsNumber() bool {
	return fieldTypeDetails[f].isNumber
}

// DefaultValue returns the default value for the field type.
func (f FieldType) DefaultValue() any {
	return fieldTypeDetails[f].defaultValue
}

// String implements the fmt.Stringer interface for easy printing.
func (f FieldType) String() string {
	return f.Name()
}

func GetFieldTypeByName(name string) (FieldType, bool) {
	return GetFieldTypeByNameIgnoreCase(name)
}

// GetFieldTypeByNameIgnoreCase finds a FieldType by its name, case-insensitively.
// It includes a special case to handle "int" as an alias for Integer.
func GetFieldTypeByNameIgnoreCase(name string) (FieldType, bool) {
	ft, ok := nameMap[strings.ToUpper(name)]
	if ok {
		return ft, true
	}

	if strings.EqualFold(name, "int") {
		return FieldTypeInteger, true
	}

	return 0, false
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
