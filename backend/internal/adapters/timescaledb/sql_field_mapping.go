package timescaledb

import (
	"backend/internal/types"
	"fmt"
	"sort"
	"strings"
)

// getFieldMappings 获取字段映射
func (g *SQLGenerator) getFieldMappings(
	unsInfo types.UnsInfo,
	viewColumns []ViewColumnInfo,
) (map[string]string, []string, map[string]map[int]bool) {

	// 构建现有映射：视图字段名 -> 物理表字段名
	existingMap := make(map[string]string)

	if len(viewColumns) > 0 {
		for _, col := range viewColumns {
			existingMap[col.ColumnName] = col.SourceColumn
		}
		for _, field := range unsInfo.GetFields() {
			if field.IsSystemField() {
				continue
			}
			fieldName := field.Name
			if sourceCol, exists := existingMap[fieldName]; exists {
				fieldType := types.FieldType(field.Type)
				prefix := fieldTypeToPrefix[fieldType]
				if !strings.HasPrefix(sourceCol, prefix) {
					delete(existingMap, fieldName) //删除类型不兼容的老的字段映射
				}
			}
		}
	}
	// 按类型统计已使用的编号
	usedNumbers := make(map[string]map[int]bool)
	for fieldType := range fieldTypeToPrefix {
		prefix := fieldTypeToPrefix[fieldType]
		usedNumbers[prefix] = g.getUsedFieldNumbers(existingMap, prefix)
	}

	// 新的字段映射
	newMappings := make(map[string]string)
	var removedFields []string

	// 为每个 Uns 字段分配物理表字段
	for _, field := range unsInfo.GetFields() {
		if field.IsSystemField() {
			continue
		}
		fieldName := field.Name

		// 如果已有映射，使用现有映射
		if sourceCol, exists := existingMap[fieldName]; exists {
			newMappings[fieldName] = sourceCol
			field.Index = &sourceCol
			delete(existingMap, fieldName) // 从 existingMap 中移除，剩余的表示被删除的字段
		} else {
			// 分配新的编号
			fieldType := types.FieldType(field.Type)
			prefix := fieldTypeToPrefix[fieldType]
			used := usedNumbers[prefix]
			num := g.allocateFieldNumber(used)
			used[num] = true

			sourceCol := fmt.Sprintf("%s_%d", prefix, num)
			newMappings[fieldName] = sourceCol

			// 设置 Index 属性
			if field.Index == nil {
				index := sourceCol
				field.Index = &index
			}
		}
	}

	// 收集被删除的字段
	for _, sourceCol := range existingMap {
		removedFields = append(removedFields, sourceCol)
		// 这里可以记录日志，但不在 SQL 生成中处理
	}

	return newMappings, removedFields, usedNumbers
}

// AnalyzeRequiredFields 分析需要的字段
func (g *SQLGenerator) AnalyzeRequiredFields(
	unsList []UnsViewInfo,
	physicsTableFields []*types.FieldDefine,
) *RequiredFields {

	fieldsToAdd := make(map[string]bool)
	fieldTypeMap := make(map[string]types.FieldType)

	for _, unsView := range unsList {
		unsInfo := unsView.Uns
		viewInfo := unsView.View

		// 获取字段映射
		mappings, _, _ := g.getFieldMappings(unsInfo, viewInfo.Columns)

		// 收集需要的字段
		for viewField, sourceCol := range mappings {
			if !g.hasPhysicalField(physicsTableFields, sourceCol) {
				if !fieldsToAdd[sourceCol] {
					fieldsToAdd[sourceCol] = true

					// 记录字段类型
					for _, field := range unsInfo.GetFields() {
						if !field.IsSystemField() && field.Name == viewField {
							fieldTypeMap[sourceCol] = types.FieldType(field.Type)
							break
						}
					}
				}
			}
		}
	}

	// 转换为切片并排序，确保生成的 SQL 顺序一致
	var fieldsSlice []string
	for field := range fieldsToAdd {
		fieldsSlice = append(fieldsSlice, field)
	}
	sort.Strings(fieldsSlice)

	return &RequiredFields{
		FieldsToAdd:  fieldsSlice,
		FieldTypeMap: fieldTypeMap,
	}
}

// 字段类型到物理表前缀的映射
var fieldTypeToPrefix = map[types.FieldType]string{
	types.FieldTypeInteger:  "long",
	types.FieldTypeLong:     "long",
	types.FieldTypeFloat:    "double",
	types.FieldTypeDouble:   "double",
	types.FieldTypeBoolean:  "bool",
	types.FieldTypeDatetime: "date",
	types.FieldTypeString:   "str",
}

// 字段类型到 PostgreSQL 数据类型的映射
var fieldTypeToPgType = map[types.FieldType]string{
	types.FieldTypeInteger:  "int4",
	types.FieldTypeLong:     "int8",
	types.FieldTypeFloat:    "float4",
	types.FieldTypeDouble:   "float8",
	types.FieldTypeBoolean:  "boolean",
	types.FieldTypeDatetime: "timestamptz",
	types.FieldTypeString:   "varchar",
}
