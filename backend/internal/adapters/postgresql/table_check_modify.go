package postgresql

import (
	"backend/internal/types"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

// 常量定义
const (
	MDF_NEW_TABLE    = 0
	MDF_TYPE_CHANGED = 1
	MDF_ADD_OR_DEL   = 2
	NO_CHANGE        = 3
)

// 检查表修改
func checkTableModify(dto types.UnsInfo, dbName string, tableName string, alterTableSQLs *[]string, tableInfo *TableInfo) int {
	if tableInfo == nil || len(tableInfo.FieldTypes) == 0 {
		return MDF_NEW_TABLE
	}

	// 构建当前字段类型映射
	curFieldTypes := make(map[string]*types.FieldDefine)
	for _, field := range dto.GetFields() {
		curFieldTypes[field.GetName()] = field
	}

	hasTypeChanged := false
	delFs := []string{}

	// 检查字段变化
	for field, oldType := range tableInfo.FieldTypes {
		curType, exists := curFieldTypes[field]
		if !exists {
			delFs = append(delFs, field)
			continue
		}

		// 从映射中移除已存在的字段
		delete(curFieldTypes, field)

		if oldType != "" {
			newTypeMid := getTypeDefineWithoutLen(curType)
			newType := GetFieldType(newTypeMid) // 使用之前定义的 GetFieldType 函数
			if oldType != newType {
				hasTypeChanged = true
				logx.Infof("typeChange %s.%s: %s->%s\n", tableName, field, oldType, newTypeMid)
				break
			}
		}
	}

	if hasTypeChanged {
		// 修改字段类型的情况则删除表
		dropSQL := fmt.Sprintf(`drop table IF EXISTS "%s"."%s"`, dbName, tableName)
		*alterTableSQLs = append(*alterTableSQLs, dropSQL)
		return MDF_TYPE_CHANGED
	} else if len(delFs) > 0 || len(curFieldTypes) > 0 {
		// pg 删除或新增字段
		alterSQL := fmt.Sprintf(`ALTER TABLE "%s"."%s"`, dbName, tableName)

		typeIds := make(map[string]map[string]bool)

		// 处理唯一约束字段
		for _, def := range curFieldTypes {
			if def.IsUnique() {
				typeName := def.GetType().Name()
				if _, exists := typeIds[typeName]; !exists {
					typeIds[typeName] = make(map[string]bool)
				}
				typeIds[typeName][def.GetName()] = true
			}
		}

		// 处理删除的字段
		for _, delF := range delFs {
			oldType := tableInfo.FieldTypes[delF]
			rename := ""

			if typeSet, exists := typeIds[oldType]; exists && len(typeSet) > 0 {
				// 获取一个可重命名的字段
				for renameField := range typeSet {
					rename = renameField
					delete(typeSet, renameField)
					delete(curFieldTypes, renameField)
					break
				}
			}

			if rename == "" {
				alterSQL += fmt.Sprintf(` DROP IF EXISTS "%s",`, delF)
			} else {
				rmSql := fmt.Sprintf(`ALTER TABLE "%s"."%s" RENAME COLUMN "%s" TO "%s"`,
					dbName, tableName, delF, rename)
				*alterTableSQLs = append(*alterTableSQLs, rmSql)
			}
		}

		// 处理新增的字段
		for _, def := range curFieldTypes {
			field := def.GetName()
			typeStr := getTypeDefine(def)
			alterSQL += fmt.Sprintf(` ADD IF NOT EXISTS "%s" %s,`, field, typeStr)
		}

		// 移除最后一个逗号并添加 SQL
		if len(alterSQL) > 0 && alterSQL[len(alterSQL)-1] == ',' {
			alterSQL = alterSQL[:len(alterSQL)-1]
		}
		*alterTableSQLs = append(*alterTableSQLs, alterSQL)
		return MDF_ADD_OR_DEL
	}

	return NO_CHANGE
}

// 获取无长度限制的类型定义
func getTypeDefineWithoutLen(def *types.FieldDefine) string {
	typeStr := getTypeDefineWithSerial(def, false)
	typeStr = strings.ToLower(typeStr)

	// 移除长度信息
	if idx := strings.Index(typeStr, "("); idx > 0 {
		typeStr = typeStr[:idx]
	}

	return typeStr
}

// 获取类型定义
func getTypeDefine(def *types.FieldDefine) string {
	return getTypeDefineWithSerial(def, true)
}

// 获取类型定义（带序列处理）
func getTypeDefineWithSerial(def *types.FieldDefine, procSerial bool) string {
	typeName := def.GetType()
	typeStr := fieldType2DBTypeMap[typeName.Name()]

	switch typeName {
	case types.FieldTypeString:
		nameLower := strings.ToLower(def.GetName())
		if strings.Contains(nameLower, "json") {
			typeStr = "jsonb"
		} else if def.GetMaxLen() != nil {
			typeStr = fmt.Sprintf("varchar(%d)", *def.GetMaxLen())
		} else {
			typeStr = "varchar(64)"
		}

	case types.FieldTypeDatetime:
		typeStr = "timestamptz(3)"

	case types.FieldTypeInteger:
		if procSerial && def.IsUnique() {
			typeStr = "serial"
		}

	case types.FieldTypeLong:
		if procSerial && def.IsUnique() && def.GetTbValueName() == nil {
			typeStr = "bigserial"
		}
	}

	return typeStr
}
