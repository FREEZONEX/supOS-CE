package postgresql

import (
	"backend/internal/types"
	"backend/share/base"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type TableInfo struct {
	_pkSet     map[string]byte
	PKs        []string          // 主键数组
	FieldTypes map[string]string // 字段名到类型的映射
}

func newTableInfo() *TableInfo {
	return &TableInfo{
		_pkSet:     make(map[string]byte),
		FieldTypes: make(map[string]string),
	}
}

type queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// ListTableInfos 列出表信息
func ListTableInfos(query queryer, topics []*types.CreateTopicDto) (map[string]*TableInfo, error) {
	if topics == nil || len(topics) == 0 {
		return make(map[string]*TableInfo), nil
	}

	// 按 schema 分组表名
	schemaTables := make(map[string]map[string]bool)
	for _, dto := range topics {
		tableName := dto.GetTable()
		dot := strings.Index(tableName, ".")

		dbName := "public" // 默认 schema，根据实际情况调整
		if dot > 0 {
			dbName = tableName[:dot]
			tableName = tableName[dot+1:]
		}

		if _, exists := schemaTables[dbName]; !exists {
			schemaTables[dbName] = make(map[string]bool)
		}
		schemaTables[dbName][tableName] = true
	}

	result := make(map[string]*TableInfo)

	// 遍历每个 schema
	for schema, tablesSet := range schemaTables {
		// 提取表名列表
		var tables []string
		for table := range tablesSet {
			tables = append(tables, table)
		}

		tableInfoMap, err := listSchemaTableInfos(query, schema, tables)
		if err != nil {
			return nil, err
		}

		// 合并结果
		for k, v := range tableInfoMap {
			result[k] = v
		}
	}

	return result, nil
}

// 查询指定 schema 和表的详细信息
func listSchemaTableInfos(query queryer, schema string, tables []string) (map[string]*TableInfo, error) {
	allMap := make(map[string]*TableInfo)

	// 分批处理，每批最多 999 个表
	batchSize := 999
	for i := 0; i < len(tables); i += batchSize {
		end := i + batchSize
		if end > len(tables) {
			end = len(tables)
		}
		batchTables := tables[i:end]

		// 查询字段信息
		err := queryColumnInfo(query, schema, batchTables, allMap)
		if err != nil {
			return nil, err
		}

		// 查询主键信息
		err = queryPrimaryKeys(query, schema, batchTables, allMap)
		if err != nil {
			return nil, err
		}

		// 更新主键数组
		for _, info := range allMap {
			info.PKs = base.MapKeys(info._pkSet)
		}
	}

	return allMap, nil
}

// 查询字段信息
func queryColumnInfo(query queryer, schema string, tables []string, allMap map[string]*TableInfo) error {
	// 构建查询语句
	var placeholders []string
	params := []interface{}{schema}

	for i, table := range tables {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+2))
		params = append(params, table)
	}

	sql := fmt.Sprintf(`
		SELECT table_name, column_name, udt_name 
		FROM information_schema.columns 
		WHERE table_name IN (%s) 
		AND table_schema = $1 
		ORDER BY table_name, column_name`,
		strings.Join(placeholders, ","))

	rows, err := query.Query(context.Background(), sql, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName, udtName string
		err := rows.Scan(&tableName, &columnName, &udtName)
		if err != nil {
			return err
		}

		if _, exists := allMap[tableName]; !exists {
			allMap[tableName] = newTableInfo()
		}

		fieldType := GetFieldType(udtName)
		allMap[tableName].FieldTypes[columnName] = fieldType
	}

	return rows.Err()
}

// 查询主键信息
func queryPrimaryKeys(query queryer, schema string, tables []string, allMap map[string]*TableInfo) error {
	// 构建查询语句
	var placeholders []string
	params := []interface{}{schema}

	for i, table := range tables {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+2))
		params = append(params, table)
	}

	sql := fmt.Sprintf(`
		SELECT 
			tc.table_name,
			kcu.column_name,
			CASE WHEN tc.constraint_type = 'PRIMARY KEY' THEN true ELSE false END AS is_primary
		FROM information_schema.table_constraints tc  
		JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name
		WHERE tc.table_name IN (%s)
		AND tc.table_schema = $1`,
		strings.Join(placeholders, ","))

	rows, err := query.Query(context.Background(), sql, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName string
		var isPrimary bool

		err := rows.Scan(&tableName, &columnName, &isPrimary)
		if err != nil {
			return err
		}

		if info, exists := allMap[tableName]; exists && isPrimary {
			info._pkSet[columnName] = 1
		}
	}

	return rows.Err()
}
