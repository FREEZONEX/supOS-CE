package timescaledb

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"context"
	"strings"
)

// 获取物理表字段（需要数据库连接）
func getPhysicalTableFields(ctx context.Context, pool queryer) (fields []*types.FieldDefine, ct, qos string, er error) {
	sql := `
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'uns_timeserial'
		ORDER BY ordinal_position;
	`

	rows, err := pool.Query(ctx, sql)
	if err != nil {
		// 表可能不存在
		return nil, "", "", nil
	}
	defer rows.Close()

	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			return fields, "", "", err
		}
		if columnName == constants.SystemSeqTag {
			continue
		}
		// 跳过系统字段

		// 转换为 FieldType
		var fieldType types.FieldType
		switch dataType {
		case "bigint", "int8":
			fieldType = types.FieldTypeLong
			if !strings.Contains(columnName, "_") {
				qos = columnName
			}
		case "integer", "int4":
			fieldType = types.FieldTypeInteger
		case "double precision", "float8":
			fieldType = types.FieldTypeDouble
		case "real", "float4":
			fieldType = types.FieldTypeFloat
		case "boolean":
			fieldType = types.FieldTypeBoolean
		case "timestamp with time zone", "timestamptz":
			fieldType = types.FieldTypeDatetime
			if !strings.Contains(columnName, "_") {
				ct = columnName
			}
		case "character varying", "varchar", "text":
			fieldType = types.FieldTypeString
		default:
			if strings.HasPrefix(dataType, "timestamptz") {
				fieldType = types.FieldTypeDatetime
			} else {
				fieldType = types.FieldTypeDouble
			}
		}

		field := &types.FieldDefine{
			Name: columnName,
			Type: fieldType.Name(),
		}
		fields = append(fields, field)
	}

	return fields, ct, qos, nil
}
