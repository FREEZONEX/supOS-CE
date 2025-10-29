package FieldFlags

// ExtendField names
var _extendFields = []string{"unit", "upperLimit", "lowerLimit", "decimal"}

// _extendFlags for extended fields
var _extendFlags = []int32{1 << 0, 1 << 1, 1 << 2, 1 << 3}

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

	for i, field := range _extendFields {
		if _, ok := used[field]; ok {
			flags |= _extendFlags[i]
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
	for i, field := range _extendFields {
		baseFlag := _extendFlags[i]
		if (*flag & baseFlag) == baseFlag {
			used = append(used, field)
		}
	}
	return used
}

// GetExtendFields returns the list of extendable field names
func GetExtendFields() []string {
	return _extendFields
}
