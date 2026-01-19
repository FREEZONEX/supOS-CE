package timescaledb

import (
	"backend/internal/types"
	"fmt"
	"sort"
	"strings"
	"time"
)

type ViewMap struct {
	viewField  string
	sourceCol  string
	expression string
}

// getFieldMappings 获取字段映射
func (g *SQLGenerator) getFieldMappings(unsInfo types.UnsInfo, viewColumns []ViewColumnInfo) []ViewMap {

	// 构建现有映射：视图字段名 -> 物理表字段名
	existingMap := make(map[string]ViewColumnInfo)
	src2colmap := make(map[string]bool)
	fieldMap := make(map[string]bool)
	if len(viewColumns) > 0 {
		for _, col := range viewColumns {
			existingMap[col.ColumnName] = col
			src2colmap[col.SourceColumn] = true
		}
		for _, field := range unsInfo.GetFields() {
			if field.IsSystemField() {
				delete(existingMap, field.Name)
				continue
			}
			fieldName := field.Name
			fieldMap[fieldName] = true
			if col, exists := existingMap[fieldName]; exists {
				fieldType := types.FieldType(field.Type)
				prefix := fieldTypeToPrefix[fieldType]
				if !strings.HasPrefix(col.SourceColumn, prefix) {
					delete(existingMap, fieldName) //删除类型不兼容的老的字段映射
				}
			}
		}
	} else {
		for _, field := range unsInfo.GetFields() {
			if !field.IsSystemField() {
				fieldMap[field.Name] = true
			}
		}
	}
	// 按类型统计已使用的编号
	viewUsedNumbers := make(map[string]map[int]bool)
	for fieldType := range fieldTypeToPrefix {
		prefix := fieldTypeToPrefix[fieldType]
		viewUsedNumbers[prefix] = g.getUsedFieldNumbers(existingMap, fieldMap, prefix)
	}

	// 新的字段映射
	newMappings := make([]ViewMap, 0, len(viewUsedNumbers))
	overrideFields := make([]string, 0, len(unsInfo.GetFields()))
	nowUnixTimestamp := time.Now().Unix()
	// 为每个 Uns 字段分配物理表字段
	for _, field := range unsInfo.GetFields() {
		if field.IsSystemField() {
			continue
		}
		fieldName := field.Name

		// 如果已有映射，使用现有映射
		if col, exists := existingMap[fieldName]; exists {
			srcCol := col.SourceColumn
			newMappings = append(newMappings, ViewMap{
				viewField:  fieldName,
				sourceCol:  srcCol,
				expression: col.Expression,
			})
			field.Index = &srcCol
			overrideFields = append(overrideFields, fieldName)
		} else {
			// 分配新的编号
			fieldType := types.FieldType(field.Type)
			prefix := fieldTypeToPrefix[fieldType]
			usedNumbers := viewUsedNumbers[prefix]
			// 从1开始查找第一个未使用的编号
			num := -1
			for i := 1; ; i++ {
				if _, has := usedNumbers[i]; !has {
					num = i
					break
				}
			}
			usedNumbers[num] = true

			sourceCol := fmt.Sprintf("%s_%d", prefix, num)
			vm := ViewMap{
				viewField: fieldName,
				sourceCol: sourceCol,
			}
			if src2colmap[sourceCol] {
				vm.expression = fmt.Sprintf(
					`CASE WHEN "timeStamp" < to_timestamp(%d) THEN NULL ELSE %s END AS %s`, nowUnixTimestamp, sourceCol, fieldName)
			}
			newMappings = append(newMappings, vm)

			// 设置 Index 属性
			if field.Index == nil {
				index := sourceCol
				field.Index = &index
			}
		}
	}
	return newMappings
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
		mappings := g.getFieldMappings(unsInfo, viewInfo.Columns)

		// 收集需要的字段
		for _, vm := range mappings {
			viewField := vm.viewField
			sourceCol := vm.sourceCol
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
var prefixToFieldType = map[string]string{
	"long":   types.FieldTypeLong,
	"double": types.FieldTypeDouble,
	"bool":   types.FieldTypeBoolean,
	"date":   types.FieldTypeDatetime,
	"str":    types.FieldTypeString,
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
