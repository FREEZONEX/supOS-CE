package timescaledb

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

type ViewInfo struct {
	SrcTable string
	Columns  []ViewColumnAnalysis
}

// ViewColumnAnalysis 视图字段分析
type ViewColumnAnalysis struct {
	ColumnName    string   `json:"column_name"`
	DataType      string   `json:"data_type"`
	IsCalculated  bool     `json:"is_calculated"`
	SourceTables  []string `json:"source_tables"`
	SourceColumns []string `json:"source_columns"`
	Expression    string   `json:"expression,omitempty"`
	Position      int      `json:"position"`
}

type queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// 基于字符串分析的视图解析器
func parseViews(pool queryer, ctx context.Context, schema string, viewNames ...string) (info map[string]ViewInfo, er error) {
	// 获取视图定义
	viewDefMap, err := getViewDefinition(pool, ctx, schema, viewNames)
	if err != nil || len(viewDefMap) == 0 {
		return nil, err
	}

	// 获取视图字段
	columnsMap, err := getViewColumns(pool, ctx, schema, viewNames)
	if err != nil {
		return nil, err
	}
	info = make(map[string]ViewInfo, len(viewDefMap))
	for name, viewDef := range viewDefMap {
		// 规范化视图定义
		normalizedDef := normalizeSQL(viewDef)

		// 提取 SELECT 字段列表
		fieldExpressions := extractFieldExpressions(normalizedDef)
		columns := columnsMap[name]
		// 分析每个字段
		var analyses []ViewColumnAnalysis
		for i, col := range columns {
			analysis := ViewColumnAnalysis{
				ColumnName: col.Name,
				DataType:   col.DataType,
				Position:   i + 1,
			}

			// 查找对应的表达式
			if i < len(fieldExpressions) {
				expr := fieldExpressions[i]
				analysis.Expression = expr.Expression
				analysis.IsCalculated = expr.IsCalculated

				// 提取源表信息
				sources := extractSourcesFromExpression(expr.Expression, normalizedDef)
				analysis.SourceTables = sources.Tables
				analysis.SourceColumns = sources.Columns
			}
			analyses = append(analyses, analysis)
		}
		srcTable := parseSrcTable(normalizedDef)
		info[name] = ViewInfo{SrcTable: srcTable, Columns: analyses}
	}

	return
}

// 提取字段表达式
func extractFieldExpressions(sql string) []FieldExpr {
	var exprs []FieldExpr

	// 提取 SELECT 和 FROM 之间的部分
	selectIndex := strings.Index(strings.ToUpper(sql), "SELECT")
	fromIndex := strings.Index(strings.ToUpper(sql), "FROM")

	if selectIndex < 0 || fromIndex < 0 || selectIndex >= fromIndex {
		return exprs
	}

	fieldStr := strings.TrimSpace(sql[selectIndex+6 : fromIndex])

	// 分割字段，注意处理逗号在括号内的情况
	exprs = splitFieldList(fieldStr)
	return exprs
}

// FieldExpr 字段表达式
type FieldExpr struct {
	Expression   string
	Alias        string
	IsCalculated bool
}

// 分割字段列表
func splitFieldList(fieldStr string) []FieldExpr {
	var exprs []FieldExpr

	// 移除换行符
	fieldStr = strings.ReplaceAll(fieldStr, "\n", " ")
	fieldStr = strings.ReplaceAll(fieldStr, "\r", " ")

	// 分割字段，考虑括号
	var current strings.Builder
	parenDepth := 0
	singleQuote := false
	doubleQuote := false

	for i := 0; i < len(fieldStr); i++ {
		ch := fieldStr[i]

		// 处理引号
		if ch == '\'' && (i == 0 || fieldStr[i-1] != '\\') {
			singleQuote = !singleQuote
		} else if ch == '"' && (i == 0 || fieldStr[i-1] != '\\') {
			doubleQuote = !doubleQuote
		}

		// 处理括号
		if !singleQuote && !doubleQuote {
			if ch == '(' {
				parenDepth++
			} else if ch == ')' {
				parenDepth--
			}
		}

		// 分割字段
		if ch == ',' && parenDepth == 0 && !singleQuote && !doubleQuote {
			expr := strings.TrimSpace(current.String())
			if expr != "" {
				exprs = append(exprs, parseFieldExpression(expr))
			}
			current.Reset()
		} else {
			current.WriteByte(ch)
		}
	}

	// 最后一个字段
	expr := strings.TrimSpace(current.String())
	if expr != "" {
		exprs = append(exprs, parseFieldExpression(expr))
	}

	return exprs
}

// 解析字段表达式
func parseFieldExpression(expr string) FieldExpr {
	field := FieldExpr{Expression: expr}

	// 检查是否有 AS 别名
	asRegex := regexp.MustCompile(`(?i)\s+AS\s+["']?([\w_]+)["']?$`)
	matches := asRegex.FindStringSubmatch(expr)
	if len(matches) > 1 {
		field.Alias = matches[1]
		field.Expression = strings.TrimSuffix(expr, matches[0])
	} else {
		// 检查是否有隐式别名（字段名本身）
		// 这里需要更复杂的逻辑，暂时跳过
	}

	// 判断是否为计算字段
	field.IsCalculated = isCalculatedExpression(field.Expression)

	return field
}

// 判断是否为计算表达式
func isCalculatedExpression(expr string) bool {
	expr = strings.ToLower(strings.TrimSpace(expr))

	// 简单标识符（表名.字段名 或 字段名）
	simpleIdentifier := regexp.MustCompile(`^[a-z_][a-z0-9_]*(\.[a-z_][a-z0-9_]*)?$`)
	if simpleIdentifier.MatchString(expr) {
		return false
	}

	// 包含以下内容的认为是计算字段
	patterns := []string{
		`\s*\+\s*`, `\s*\-\s*`, `\s*\*\s*`, `\s*/\s*`, // 算术运算
		`\|\|`,                   // 字符串连接
		`case\s+when`,            // CASE 语句
		`coalesce\(`, `nullif\(`, // 函数
		`count\(`, `sum\(`, `avg\(`, `min\(`, `max\(`, // 聚合函数
		`extract\(`, `date_part\(`, // 时间函数
		`cast\(`, `::`, // 类型转换
	}

	for _, pattern := range patterns {
		if regexp.MustCompile(pattern).MatchString(expr) {
			return true
		}
	}

	return false
}

// ExpressionSources 表达式来源
type ExpressionSources struct {
	Tables  []string
	Columns []string
}

// 从表达式中提取源信息
func extractSourcesFromExpression(expr, viewDef string) ExpressionSources {
	sources := ExpressionSources{}

	// 提取表名.字段名模式
	pattern := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\.([a-zA-Z_][a-zA-Z0-9_]*)`)
	matches := pattern.FindAllStringSubmatch(expr, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			if len(match) >= 3 {
				tableName := match[1]
				columnName := match[2]

				// 验证表名是否在 FROM 子句中
				if isTableInFromClause(tableName, viewDef) {
					// 添加表（去重）
					found := false
					for _, t := range sources.Tables {
						if t == tableName {
							found = true
							break
						}
					}
					if !found {
						sources.Tables = append(sources.Tables, tableName)
					}

					// 添加字段
					sources.Columns = append(sources.Columns, columnName)
				}
			}
		}

	} else {
		columnName := ""
		if expr[0] == '(' {
			end := strings.IndexByte(expr, ')')
			columnName = expr[1:end]
		} else {
			matches = validNamePattern.FindAllStringSubmatch(expr, -1)
			if len(matches) > 0 {
				columnName = matches[0][0]
			}
		}
		sources.Columns = append(sources.Columns, columnName)
	}

	return sources
}

var validNamePattern = regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)`)

// 检查表名是否在 FROM 子句中
func isTableInFromClause(tableName, viewDef string) bool {
	// 转换为大写查找
	upperDef := strings.ToUpper(viewDef)
	upperTable := strings.ToUpper(tableName)

	// 查找 FROM 子句
	fromIndex := strings.Index(upperDef, "FROM")
	if fromIndex < 0 {
		return false
	}

	// 提取 FROM 之后的部分（直到 WHERE、GROUP BY 等）
	fromSection := viewDef[fromIndex+4:]

	// 查找下一个关键字
	keywords := []string{"WHERE", "GROUP BY", "HAVING", "ORDER BY", "LIMIT"}
	endIndex := len(fromSection)
	for _, keyword := range keywords {
		if idx := strings.Index(strings.ToUpper(fromSection), keyword); idx >= 0 && idx < endIndex {
			endIndex = idx
		}
	}

	fromSection = fromSection[:endIndex]

	// 在 FROM 部分中查找表名
	return strings.Contains(strings.ToUpper(fromSection), upperTable)
}

// 规范化 SQL
func normalizeSQL(sql string) string {
	// 移除多余空格和换行
	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\r", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")

	// 合并多个空格
	spaceRegex := regexp.MustCompile(`\s+`)
	sql = spaceRegex.ReplaceAllString(sql, " ")

	return strings.TrimSpace(sql)
}

// 数据库查询方法（与之前相同）
func getViewDefinition(pool queryer, ctx context.Context, schema string, viewNames []string) (viewDefMap map[string]string, err error) {
	query := `select viewname,definition FROM pg_views WHERE schemaname=$1 AND viewname IN ` + `('` + strings.Join(viewNames, `','`) + `')`
	var viewName, definition string
	viewDefMap = make(map[string]string, len(viewNames))
	rows, err := pool.Query(ctx, query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		er := rows.Scan(&viewName, &definition)
		if er != nil {
			return nil, er
		}
		viewDefMap[viewName] = definition
	}
	return viewDefMap, err
}

type ColumnInfo struct {
	Name     string
	DataType string
}

func getViewColumns(pool queryer, ctx context.Context, schema string, viewNames []string) (colMap map[string][]ColumnInfo, err error) {
	query := `
		SELECT 
			table_name,column_name,data_type
		FROM information_schema.columns
		WHERE table_schema = $1
			AND table_name IN %s
		ORDER BY ordinal_position
	`
	vs := `('` + strings.Join(viewNames, `','`) + `')`
	rows, err := pool.Query(ctx, fmt.Sprintf(query, vs), schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	colMap = make(map[string][]ColumnInfo, len(viewNames))
	for rows.Next() {
		var tableName string
		var col ColumnInfo
		if err := rows.Scan(&tableName, &col.Name, &col.DataType); err != nil {
			return nil, err
		}
		columns, has := colMap[tableName]
		if !has {
			columns = make([]ColumnInfo, 0, 8)
		}
		columns = append(columns, col)
		colMap[tableName] = columns
	}

	return colMap, nil
}
func parseSrcTable(viewDef string) string {
	view := strings.ToUpper(viewDef)
	fromI := strings.Index(view, "FROM")
	if fromI > 0 {
		start := fromI + 4
		first := true
		for i, v := range view[start:] {
			if unicode.IsSpace(v) {
				if first {
					continue
				} else {
					return strings.TrimSpace(viewDef[start : start+i])
				}
			} else {
				first = false
			}
		}
	}
	return ""
}
