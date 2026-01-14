package timescaledb

import (
	"context"
	"strings"
)

// 解析UNS视图（需要数据库连接）
func parseUnsViews(pool queryer, ctx context.Context, schema string, viewNames []string) (info map[string]SimpleViewInfo, er error) {
	vs, er := parseViews(pool, ctx, schema, viewNames...)
	if er != nil || len(vs) == 0 {
		return nil, er
	}
	info = make(map[string]SimpleViewInfo, len(vs))
	for name, v := range vs {
		cols := make([]ViewColumnInfo, 0, len(v.Columns))
		for _, col := range v.Columns {
			if len(col.SourceColumns) > 0 {
				expr := col.Expression
				if strings.HasPrefix(expr, `"`) && strings.HasSuffix(expr, `"`) {
					expr = expr[1 : len(expr)-1]
				}
				srcCol := col.SourceColumns[0]
				if expr == srcCol {
					expr = ""
				}
				cols = append(cols, ViewColumnInfo{
					ColumnName:   col.ColumnName,
					SourceColumn: srcCol,
					Expression:   expr,
				})
			}
		}
		info[name] = SimpleViewInfo{SrcTable: v.SrcTable, Columns: cols}
	}
	return
}
