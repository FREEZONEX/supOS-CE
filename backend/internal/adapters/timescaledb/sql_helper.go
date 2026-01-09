package timescaledb

import (
	"backend/internal/types"
	"fmt"
	"strings"
)

// getUsedFieldNumbers 获取已使用的字段编号
func (g *SQLGenerator) getUsedFieldNumbers(
	existingMap map[string]string,
	prefix string,
) map[int]bool {
	used := make(map[int]bool)
	if len(existingMap) > 0 {
		for _, src := range existingMap {
			if strings.HasPrefix(src, prefix+"_") {
				// 提取编号
				var num int
				fmt.Sscanf(src, prefix+"_%d", &num)
				if num > 0 {
					used[num] = true
				}
			}
		}
	}
	return used
}

// allocateFieldNumber 分配新的字段编号
func (g *SQLGenerator) allocateFieldNumber(usedNumbers map[int]bool) int {
	// 从1开始查找第一个未使用的编号
	for i := 1; ; i++ {
		if _, has := usedNumbers[i]; !has {
			return i
		}
	}
}

// hasPhysicalField 检查物理表是否已有某个字段
func (g *SQLGenerator) hasPhysicalField(
	physicsTableFields []*types.FieldDefine,
	fieldName string,
) bool {
	for _, field := range physicsTableFields {
		if field.Name == fieldName {
			return true
		}
	}
	return false
}

// getPgTypeForField 获取字段的 PostgreSQL 类型
func (g *SQLGenerator) getPgTypeForField(
	fieldName string,
	fieldType types.FieldType,
) string {
	// 根据前缀判断
	if strings.HasPrefix(fieldName, "long_") {
		return "int8"
	} else if strings.HasPrefix(fieldName, "double_") {
		return "float8"
	} else if strings.HasPrefix(fieldName, "bool_") {
		return "boolean"
	} else if strings.HasPrefix(fieldName, "date_") {
		return "timestamptz"
	} else if strings.HasPrefix(fieldName, "str_") {
		return "text"
	}

	// 根据字段类型判断
	if pgType, exists := fieldTypeToPgType[fieldType]; exists {
		return pgType
	}

	// 默认类型
	return "float8"
}
