package FieldUtils

import (
	"backend/internal/common"
	"backend/internal/common/constants"
	"backend/internal/common/dto"
	"backend/internal/common/enums"
	"backend/internal/common/utils/PathUtil"
	"fmt"
	"strings"
)

// ExtendField names
var ExtendFields = []string{"unit", "upperLimit", "lowerLimit", "decimal"}

// ExtendFlags for extended fields
var ExtendFlags = []int32{1 << 0, 1 << 1, 1 << 2, 1 << 3}

// GetTimestampField finds the timestamp field in field definitions
func GetTimestampField(fields []*dto.FieldDefine) *dto.FieldDefine {
	for _, f := range fields {
		if f.Type == enums.FieldTypeDatetime {
			return f
		}
	}
	return nil
}

// GetQualityField finds the quality field in field definitions
func GetQualityField(fields []*dto.FieldDefine, dataType int16) *dto.FieldDefine {
	if dataType == constants.TimeSequenceType && fields != nil && len(fields) > 2 {
		// Quality field is usually the last or second-to-last field
		lastIdx := len(fields) - 1
		if fields[lastIdx].Type == enums.FieldTypeLong {
			return fields[lastIdx]
		}
		if lastIdx > 0 && fields[lastIdx-1].Type == enums.FieldTypeLong {
			return fields[lastIdx-1]
		}
	}
	return nil
}

// ValidateFields validates a slice of field definitions
func ValidateFields(fields []*dto.FieldDefine, checkSysField bool) error {
	seen := make(map[string]bool)
	for _, f := range fields {
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
			if checkSysField {
				if _, ok := constants.SystemFields[name]; !ok {
					return fmt.Errorf("field name '%s' starts with '_' but is not a valid system field", name)
				}
			}
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
	if createTimeField != nil && createTimeField.Type != enums.FieldTypeDatetime {
		return fmt.Errorf("field '%s' must be of type DATETIME", constants.SysFieldCreateTime)
	}

	return nil
}

// CountNumericFields counts the number of numeric fields
func CountNumericFields(fields []*dto.FieldDefine) int {
	total := 0
	for _, f := range fields {
		if f.Type.IsNumber() {
			total++
		}
	}
	return total
}

// GenerateFlag generates a bitmask flag from a list of used field names
func GenerateFlag(extendFieldUsed []string) int32 {
	flags := int32(0)
	if len(extendFieldUsed) == 0 {
		return flags
	}
	used := make(map[string]struct{}, len(extendFieldUsed))
	for _, field := range extendFieldUsed {
		used[field] = struct{}{}
	}

	for i, field := range ExtendFields {
		if _, ok := used[field]; ok {
			flags |= ExtendFlags[i]
		}
	}
	return flags
}

// ParseFlag parses a bitmask flag into a list of used field names
func ParseFlag(flag *int32) []string {
	if flag == nil || *flag == 0 {
		return nil
	}
	var used []string
	for i, field := range ExtendFields {
		baseFlag := ExtendFlags[i]
		if (*flag & baseFlag) == baseFlag {
			used = append(used, field)
		}
	}
	return used
}

// GetExtendFields returns the list of extendable field names
func GetExtendFields() []string {
	return ExtendFields
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
func SetDefaultMaxLen(field *dto.FieldDefine) {
	if field.Type == enums.FieldTypeString && field.MaxLen == 0 {
		nameLower := strings.ToLower(field.Name)
		defaultLen := dto.DefaultMaxStrLen

		// Set shorter default for name/tag fields
		if strings.Contains(nameLower, "name") || strings.Contains(nameLower, "tag") {
			defaultLen = 64
		}

		field.MaxLen = defaultLen
	}
}

// CountNumberFields counts the number of numeric fields
func CountNumberFields(fields []*dto.FieldDefine) int {
	count := 0
	for _, f := range fields {
		if f.Type.IsNumber() && !f.IsSystemField() {
			count++
		}
	}
	return count
}

// FilterBlobFields filters all BLOB and LBLOB fields
func FilterBlobFields(fields []*dto.FieldDefine) []*dto.FieldDefine {
	result := make([]*dto.FieldDefine, 0)
	for _, f := range fields {
		if f.Type == enums.FieldTypeBlob || f.Type == enums.FieldTypeLBlob {
			result = append(result, f)
		}
	}
	return result
}

// FindFieldByName finds field by name
func FindFieldByName(fields []*dto.FieldDefine, name string) *dto.FieldDefine {
	for _, f := range fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// HasDuplicateFieldNames checks for duplicate field names
func HasDuplicateFieldNames(fields []*dto.FieldDefine) (string, bool) {
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
	Fields    []*dto.FieldDefine
}

// ProcessFieldDefines validates and processes a list of field definitions, optionally adding system fields.
func ProcessFieldDefines(jdbcType common.SrcJdbcType, fields []*dto.FieldDefine, checkSysField bool, addSysField bool) (*TableFieldDefine, error) {
	if len(fields) == 0 {
		return nil, nil
	}

	// Create a deep copy to avoid modifying the original slice and its elements
	processedFields := make([]*dto.FieldDefine, len(fields))
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
	var fNews []*dto.FieldDefine

	if jdbcType.TypeCode() == constants.TimeSequenceType {
		// Time-series data
		fNews = make([]*dto.FieldDefine, 0, len(processedFields)+2)
		fNews = append(fNews, &dto.FieldDefine{Name: constants.SysFieldCreateTime, Type: enums.FieldTypeDatetime})

		var nonSysFields []*dto.FieldDefine
		for _, f := range processedFields {
			if !f.IsSystemField() {
				nonSysFields = append(nonSysFields, f)
			}
		}

		// Special handling for TimeScaleDB when there is only one non-system field named "value".
		if jdbcType == common.SrcJdbcTypeTimeScaleDB && len(nonSysFields) == 1 && nonSysFields[0].Name == constants.SystemSeqValue {
			theField := nonSysFields[0]
			theField.Name = constants.SystemSeqValue // Ensure the name is correct
			fNews = append(fNews, theField)

			tableValueField := &dto.FieldDefine{
				Name:        constants.SystemSeqTag,
				Type:        enums.FieldTypeLong,
				Unique:      true,
				TbValueName: theField.Name,
				MaxLen:      200,
			}
			fNews = append(fNews, tableValueField)

			tableName = "supos_timeserial_" + strings.ToLower(theField.Type.String())
		} else {
			// Default behavior for other time-series data
			tableName = ""
			fNews = append(fNews, nonSysFields...)
		}

		fNews = append(fNews, &dto.FieldDefine{Name: constants.QosField, Type: enums.FieldTypeLong})
	} else {
		// Relational data
		tableName = ""
		fNews = make([]*dto.FieldDefine, 0, len(processedFields)+2)
		hasId := false

		fNews = append(fNews, &dto.FieldDefine{Name: constants.SysFieldCreateTime, Type: enums.FieldTypeDatetime})

		for _, f := range processedFields {
			if !f.IsSystemField() {
				if !hasId && f.IsUnique() {
					hasId = true
				}
				fNews = append(fNews, f)
			}
		}

		// If no unique field is defined by the user, add a system ID field.
		if !hasId {
			fNews = append(fNews, &dto.FieldDefine{Name: constants.SysFieldID, Type: enums.FieldTypeLong, Unique: true})
		}
	}

	return &TableFieldDefine{
		TableName: tableName,
		Fields:    fNews,
	}, nil
}
