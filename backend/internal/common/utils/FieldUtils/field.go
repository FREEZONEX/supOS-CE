package FieldUtils

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/types"
	"errors"
	"fmt"
	"strings"
)

// GetTimestampField finds the timestamp field in field definitions
func GetTimestampField(fields []*types.FieldDefine) *types.FieldDefine {
	for _, f := range fields {
		if f.Type == types.FieldTypeDatetime {
			return f
		}
	}
	return nil
}

// GetQualityField finds the quality field in field definitions
func GetQualityField(fields []*types.FieldDefine, dataType int16) *types.FieldDefine {
	if dataType == constants.TimeSequenceType && fields != nil && len(fields) > 2 {
		// Quality field is usually the last or second-to-last field
		lastIdx := len(fields) - 1
		if fields[lastIdx].Type == types.FieldTypeLong {
			return fields[lastIdx]
		}
		if lastIdx > 0 && fields[lastIdx-1].Type == types.FieldTypeLong {
			return fields[lastIdx-1]
		}
	}
	return nil
}

// ValidateFields validates a slice of field definitions
func ValidateFields(fields []*types.FieldDefine, checkSysField bool) error {
	seen := make(map[string]bool)
	for _, f := range fields {
		if fieldType, ok := types.GetFieldTypeByName(f.Type); ok {
			f.Type = fieldType.Name()
		} else {
			return errors.New(I18nUtils.GetMessage("uns.invalid.type", f.Type))
		}
		name := strings.TrimSpace(f.Name)
		if name == "" {
			return fmt.Errorf("field name cannot be empty")
		}
		if len(name) > 63 {
			return fmt.Errorf("field name '%s' is too long (max 63)", name)
		}
		if seen[name] {
			return fmt.Errorf("duplicate field name '%s'", name)
		}
		seen[name] = true
		f.Name = name

		if f.IsSystemField() {
			continue
		}

		if name[0] >= '0' && name[0] <= '9' {
			return fmt.Errorf("field name '%s' cannot start with a digit", name)
		}
		if !PathUtil.IsFieldNameFormatOK(name) {
			return fmt.Errorf("field name '%s' has invalid format", name)
		}
	}

	createTimeField := FindFieldByName(fields, constants.SysFieldCreateTime)
	if createTimeField != nil && createTimeField.Type != types.FieldTypeDatetime {
		return fmt.Errorf("field '%s' must be of type DATETIME", constants.SysFieldCreateTime)
	}

	return nil
}

// CountNumericFields counts the number of numeric fields
func CountNumericFields(fields []*types.FieldDefine) int {
	total := 0
	for _, f := range fields {
		if types.FieldType(f.Type).IsNumber() {
			total++
		}
	}
	return total
}

// ValidateFieldName is deprecated, use ValidateFields instead
func ValidateFieldName(name string) error {
	if name == "" {
		return fmt.Errorf("field name cannot be empty")
	}

	if len(name) > 63 {
		return fmt.Errorf("field name '%s' is too long (max 63 characters)", name)
	}

	// Check if starts with system field prefix
	if strings.HasPrefix(name, constants.SystemFieldPrev) {
		if _, ok := constants.SystemFields[name]; !ok {
			return fmt.Errorf("field name '%s' cannot start with '_' unless it's a system field", name)
		}
	}

	// Check if starts with digit
	if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
		return fmt.Errorf("field name '%s' cannot start with a digit", name)
	}

	return nil
}

// SetDefaultMaxLen sets default max length for string fields
func SetDefaultMaxLen(field *types.FieldDefine) {
	if field.Type == types.FieldTypeString && field.MaxLen == nil {
		nameLower := strings.ToLower(field.Name)
		defaultLen := types.DefaultMaxStrLen

		// Set shorter default for name/tag fields
		if strings.Contains(nameLower, "name") || strings.Contains(nameLower, "tag") {
			defaultLen = 64
		}

		field.MaxLen = &defaultLen
	}
}

// CountNumberFields counts the number of numeric fields
func CountNumberFields(fields []*types.FieldDefine) int {
	count := 0
	for _, f := range fields {
		if types.FieldType(f.Type).IsNumber() && !f.IsSystemField() {
			count++
		}
	}
	return count
}

// FilterBlobFields filters all BLOB and LBLOB fields
func FilterBlobFields(fields []*types.FieldDefine) []*types.FieldDefine {
	result := make([]*types.FieldDefine, 0)
	for _, f := range fields {
		if f.Type == types.FieldTypeBlob || f.Type == types.FieldTypeLBlob {
			result = append(result, f)
		}
	}
	return result
}

// FindFieldByName finds field by name
func FindFieldByName(fields []*types.FieldDefine, name string) *types.FieldDefine {
	for _, f := range fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// HasDuplicateFieldNames checks for duplicate field names
func HasDuplicateFieldNames(fields []*types.FieldDefine) (string, bool) {
	seen := make(map[string]bool)
	for _, f := range fields {
		if seen[f.Name] {
			return f.Name, true
		}
		seen[f.Name] = true
	}
	return "", false
}

// TableFieldDefine holds a table name and its field definitions.
type TableFieldDefine struct {
	TableName string
	Fields    []*types.FieldDefine
}

var _True = true
var _False = false

// ProcessFieldDefines validates and processes a list of field definitions, optionally adding system fields.
func ProcessFieldDefines(jdbcType types.SrcJdbcType, fields []*types.FieldDefine, checkSysField bool, addSysField bool) (*TableFieldDefine, error) {
	if len(fields) == 0 {
		return nil, nil
	}

	// Create a deep copy to avoid modifying the original slice and its elements
	processedFields := make([]*types.FieldDefine, len(fields))
	for i, f := range fields {
		clone := *f
		processedFields[i] = &clone
	}

	// Set defaults and perform validation
	for _, f := range processedFields {
		SetDefaultMaxLen(f)
	}

	if err := ValidateFields(processedFields, checkSysField); err != nil {
		return nil, err
	}

	if !addSysField {
		return &TableFieldDefine{
			TableName: "",
			Fields:    processedFields,
		}, nil
	}

	var tableName string
	var fNews []*types.FieldDefine

	if jdbcType.TypeCode() == constants.TimeSequenceType {
		// Time-series data
		fNews = make([]*types.FieldDefine, 0, len(processedFields)+2)
		fNews = append(fNews, &types.FieldDefine{Name: constants.SysFieldCreateTime, Type: types.FieldTypeDatetime, Unique: &_True})

		var nonSysFields []*types.FieldDefine
		for _, f := range processedFields {
			if !f.IsSystemField() {
				nonSysFields = append(nonSysFields, f)
			}
		}

		// Special handling for TimeScaleDB when there is only one non-system field named "value".
		if jdbcType == types.SrcJdbcTypeTimeScaleDB {
			tableName = "uns_timeserial"
			name := constants.SystemSeqTag // Ensure the name is correct
			tableValueField := &types.FieldDefine{
				Name:        constants.SystemSeqTag,
				Type:        types.FieldTypeLong,
				Unique:      &_True,
				TbValueName: &name,
			}
			fNews = append(fNews, tableValueField)
			if len(nonSysFields) > 0 {
				fNews = append(fNews, nonSysFields...)
			}
		} else {
			// Default behavior for other time-series data
			tableName = ""
			fNews = append(fNews, nonSysFields...)
		}

		fNews = append(fNews, &types.FieldDefine{Name: constants.QosField, Type: types.FieldTypeLong})
	} else if jdbcType.TypeCode() == constants.RelationType {
		// Relational data
		tableName = ""
		fNews = make([]*types.FieldDefine, 0, len(processedFields)+2)
		hasId := false

		fNews = append(fNews, &types.FieldDefine{Name: constants.SysFieldCreateTime, Type: types.FieldTypeDatetime})

		for _, f := range processedFields {
			if !f.IsSystemField() {
				if !hasId && f.IsUnique() {
					hasId = true
				}
				fNews = append(fNews, f)
			}
		}

		// If no unique field is defined by the user, add a system Id field.
		if !hasId {
			fNews = append(fNews, &types.FieldDefine{Name: constants.SysFieldID, Type: types.FieldTypeLong, Unique: &_True})
		}
	}

	return &TableFieldDefine{
		TableName: tableName,
		Fields:    fNews,
	}, nil
}
