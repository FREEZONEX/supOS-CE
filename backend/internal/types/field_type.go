package types

import (
	"backend/share/base"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
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

var fieldTypes = []FieldType{
	FieldTypeInteger, FieldTypeLong, FieldTypeFloat, FieldTypeDouble,
	FieldTypeBoolean, FieldTypeDatetime,
	FieldTypeString, FieldTypeBlob, FieldTypeLBlob,
}

// fieldTypeInfo holds the metadata for each FieldType constant.
type fieldTypeInfo struct {
	isNumber     bool
	defaultValue any
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
	switch f {
	case FieldTypeInteger:
		return 0
	case FieldTypeLong:
		return int64(0)
	case FieldTypeFloat:
		return float32(0)
	case FieldTypeDouble:
		return float64(0)
	case FieldTypeBoolean:
		return false
	case FieldTypeDatetime:
		return time.Now()
	case FieldTypeString:
		return ""
	}
	return nil
}
func (f FieldType) ZeroValue() any {
	switch f {
	case FieldTypeInteger:
		return 0
	case FieldTypeLong:
		return int64(0)
	case FieldTypeFloat:
		return float32(0)
	case FieldTypeDouble:
		return float64(0)
	case FieldTypeBoolean:
		return false
	case FieldTypeDatetime:
		return time.Time{}
	case FieldTypeString:
		return ""
	}
	return nil
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

var sortedFieldTypes []FieldType

type fieldTypeSlice []FieldType

func (x fieldTypeSlice) Len() int           { return len(x) }
func (x fieldTypeSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x fieldTypeSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

var fieldTypeOrdinals = make(map[string]int, 16)

func init() {
	sortedFieldTypes = make([]FieldType, len(fieldTypes))
	for i, v := range fieldTypes {
		sortedFieldTypes[i] = v
		fieldTypeOrdinals[v.Name()] = i
	}
	sort.Sort(fieldTypeSlice(sortedFieldTypes))
}
func (f FieldType) GetOrdinal() int {
	return fieldTypeOrdinals[f.Name()]
}

// GetFieldTypeByNameIgnoreCase finds a FieldType by its name, case-insensitively.
// It includes a special case to handle "int" as an alias for Integer.
func GetFieldTypeByNameIgnoreCase(name string) (rs FieldType, ok bool) {
	name = strings.ToUpper(name)
	if name == "INT" {
		return FieldTypeInteger, true
	}
	k := FieldType(name)
	i := base.BinarySearchCmp(sortedFieldTypes, k)
	return k, i >= 0
}

// MarshalJSON implements the json.Marshaler interface, serializing the FieldType to its string name.
func (f FieldType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + string(f) + "\""), nil
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
