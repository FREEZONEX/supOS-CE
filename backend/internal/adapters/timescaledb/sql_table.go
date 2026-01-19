package timescaledb

import (
	"backend/internal/common/constants"
	"fmt"
	"sort"
	"strings"
)

// GenerateCreateTableSQL 生成创建表的 SQL
func (g *SQLGenerator) GenerateCreateTableSQL(
	required *RequiredFields,
) []string {

	// 基础字段
	columns := []string{
		fmt.Sprintf(`"%s" int8 NOT NULL`, constants.SystemSeqTag),
		fmt.Sprintf(`"%s" timestamptz(3) NOT NULL DEFAULT now()`, constants.SysFieldCreateTime),
		fmt.Sprintf(`"%s" int8 DEFAULT 0`, constants.QosField),
	}

	// 添加动态字段
	if required != nil {
		for _, fieldName := range required.FieldsToAdd {
			fieldType := required.FieldTypeMap[fieldName]
			pgType := g.getPgTypeForField(fieldName, fieldType)
			columns = append(columns, fmt.Sprintf(`"%s" %s NULL`, fieldName, pgType))
		}
	}

	// 确保有 double_1 字段
	hasDouble1 := false
	for _, col := range columns {
		if strings.Contains(col, `"double_1"`) {
			hasDouble1 = true
			break
		}
	}
	if !hasDouble1 {
		columns = append(columns, `"double_1" float8 NULL`)
	}

	// 构建 CREATE TABLE 语句
	columnsSQL := strings.Join(columns, ",\n    ")
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS public."uns_timeserial" (
			%s,
			PRIMARY KEY ("%s", "%s")
		);
	`, columnsSQL, constants.SystemSeqTag, constants.SysFieldCreateTime)

	return []string{sql,
		`SELECT create_hypertable (
    'uns_timeserial', 
    'timeStamp',
    partitioning_column => 'tag',
    number_partitions => 50,  
    chunk_time_interval => INTERVAL '2 hour'
    );`,

		`ALTER TABLE uns_timeserial SET (
    timescaledb.compress,                    
    timescaledb.compress_segmentby = 'tag', 
    timescaledb.compress_orderby = '"timeStamp" DESC'
);`,

		`SELECT add_retention_policy('uns_timeserial', INTERVAL '2 year'); `,

		`SELECT add_compression_policy(
    'uns_timeserial',
    compress_after => INTERVAL '1 hour',
    schedule_interval => INTERVAL '2 hour' 
);`}
}

// GenerateAlterTableSQL 生成修改表的 SQL（添加字段）
func (g *SQLGenerator) GenerateAlterTableSQL(
	required *RequiredFields,
) []string {

	if required == nil || len(required.FieldsToAdd) == 0 {
		return []string{}
	}

	// 按类型分组
	typeGroups := make(map[string][]string)

	for _, fieldName := range required.FieldsToAdd {
		fieldType := required.FieldTypeMap[fieldName]
		pgType := g.getPgTypeForField(fieldName, fieldType)

		if _, exists := typeGroups[pgType]; !exists {
			typeGroups[pgType] = []string{}
		}
		typeGroups[pgType] = append(typeGroups[pgType], fieldName)
	}

	// 为每种类型生成 SQL
	var sqls []string

	// 按键排序，确保生成的 SQL 顺序一致
	var pgTypes []string
	for pgType := range typeGroups {
		pgTypes = append(pgTypes, pgType)
	}
	sort.Strings(pgTypes)

	for _, pgType := range pgTypes {
		fields := typeGroups[pgType]
		sort.Strings(fields) // 字段名也排序

		var addClauses []string
		for _, fieldName := range fields {
			addClauses = append(addClauses,
				fmt.Sprintf(`ADD COLUMN IF NOT EXISTS "%s" %s NULL`, fieldName, pgType))
		}

		sql := `ALTER TABLE public."uns_timeserial" ` + strings.Join(addClauses, ",\n")

		sqls = append(sqls, sql)
	}

	return sqls
}
