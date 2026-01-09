package postgresql

import (
	"backend/internal/common"
	"backend/internal/common/constants"
	"backend/internal/types"
	"backend/share/base"
	"strings"
)

// 获取创建表的 SQL -- by uns
func getCreateTableSQLByUns(dto types.UnsInfo) string {
	return getCreateTableSQL(dto.GetTbFieldName() != "", dto.GetTable(), dto.GetFields())
}

// 获取创建表的 SQL
func getCreateTableSQL(isShareTable bool, tableName string, fields []*types.FieldDefine) string {
	var builder base.StringBuilder
	builder.Grow(128)
	fullName := tableName
	if tableName[len(tableName)-1] != '"' {
		fullName = GetFullTableName(tableName)
	}
	builder.Append("create table IF NOT EXISTS ").Append(fullName).Append(" (")

	ids := make(map[string]bool)

	for _, def := range fields {
		name := def.GetName()
		var typeStr string
		if def.GetType() == types.FieldTypeString && isShareTable {
			typeStr = "text"
		} else {
			typeStr = getTypeDefine(def)
		}

		if def.IsUnique() {
			ids[name] = true
		}

		builder.Append("\"").Append(name).Append("\" ").Append(typeStr)

		if def.IsUnique() {
			builder.Append(" NOT NULL ")
		} else if strings.HasPrefix(typeStr, "timestamp") && name == constants.SysSaveTime {
			builder.Append(" DEFAULT now() ")
		} else if strings.HasPrefix(typeStr, "timestamp") && name == constants.SysFieldCreateTime {
			builder.Append(" DEFAULT now() ")
		}

		builder.Append(",")
	}

	// 处理主键约束
	if len(ids) > 0 {
		table := GetCleanTableName(tableName)

		builder.Append("CONSTRAINT \"pk_").Append(table).Append("_").Long(common.NextId()).Append("\" PRIMARY KEY (")

		for pk := range ids {
			builder.Append("\"").Append(pk).Append("\",")
		}
		builder.SetLast(')').Append(" ")
	}
	builder.SetLast(')')
	sqlStr := builder.String()
	return sqlStr
}
