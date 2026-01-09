package timescaledb

import (
	"backend/internal/types"
	"fmt"
	"strings"
)

// GenerateViewSQL 生成创建视图的 SQL
func (g *SQLGenerator) GenerateViewSQL(
	unsInfo types.UnsInfo,
	mappings map[string]string,
) string {

	alias := unsInfo.GetAlias()
	id := unsInfo.GetId()

	// 构建 SELECT 字段列表
	var selectFields []string

	// 固定字段
	tbField := unsInfo.GetTbFieldName()
	for _, field := range unsInfo.GetFields() {
		if field.IsSystemField() && field.Name != tbField {
			selectFields = append(selectFields, fmt.Sprintf(`"%s"`, field.Name))
		}
	}

	// 动态字段
	for _, field := range unsInfo.GetFields() {
		sourceCol := mappings[field.Name]
		if len(sourceCol) == 0 {
			continue
		}

		// 根据字段类型可能需要类型转换
		selectExpr := fmt.Sprintf(`"%s"`, sourceCol)

		// INTEGER 类型需要转换为 int4
		if field.Type == types.FieldTypeInteger {
			selectExpr = fmt.Sprintf(`"%s"::int4`, sourceCol)
		} else if field.Type == types.FieldTypeFloat {
			// FLOAT 类型需要转换为 float4
			selectExpr = fmt.Sprintf(`"%s"::real`, sourceCol)
		}

		// 添加别名
		selectExpr = fmt.Sprintf(`%s as "%s"`, selectExpr, field.Name)
		selectFields = append(selectFields, selectExpr)
	}

	// 构建完整的 SELECT 语句
	selectClause := strings.Join(selectFields, ", ")

	// 创建视图 SQL
	sql := fmt.Sprintf(`
		CREATE OR REPLACE VIEW "%s" AS
		SELECT %s
		FROM %s
		WHERE "tag" = %d;
	`, alias, selectClause, unsInfo.GetTable(), id)

	return sql
}

// GenerateDataUpdateSQL 生成更新数据的 SQL（将删除的字段设为 NULL）
func (g *SQLGenerator) GenerateDataUpdateSQL(
	unsInfo types.UnsInfo,
	removedFields []string,
) string {

	if len(removedFields) == 0 {
		return ""
	}

	// 构建 SET 子句
	var setClauses []string
	for _, field := range removedFields {
		setClauses = append(setClauses, fmt.Sprintf(`"%s" = NULL`, field))
	}

	setClause := strings.Join(setClauses, ", ")

	sql := fmt.Sprintf(`
		UPDATE %s
		SET %s
		WHERE "tag" = %d;
	`, unsInfo.GetTable(), setClause, unsInfo.GetId())

	return sql
}
