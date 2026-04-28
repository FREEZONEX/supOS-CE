package postgresql

import (
	"backend/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultHistoryQueryLimit = 1000
	maxHistoryQueryLimit     = 10000
)

func QueryHistoryWithTable(ctx context.Context, dbPool *pgxpool.Pool, queryTable string, uns types.UnsInfo, req *types.HistoryValueRequest) (*types.UnsHistoryFileResult, error) {
	if dbPool == nil {
		return nil, fmt.Errorf("history query db pool is nil")
	}
	if uns == nil {
		return nil, fmt.Errorf("history query uns is nil")
	}
	timestampField := strings.TrimSpace(uns.GetTimestampField())
	if timestampField == "" {
		return nil, fmt.Errorf("timestamp field is empty")
	}
	if queryTable == "" {
		queryTable = uns.GetAlias()
	}

	availableColumns, err := loadHistoryColumns(ctx, dbPool, queryTable)
	if err != nil {
		return nil, err
	}
	selectColumns, selectedFields := buildHistorySelectColumns(uns, availableColumns)
	if len(availableColumns) > 0 {
		if _, ok := availableColumns[timestampField]; !ok {
			return nil, fmt.Errorf("history timestamp column not found: %s", timestampField)
		}
	}
	sql := strings.Builder{}
	sql.Grow(256)
	sql.WriteString("SELECT ")
	sql.WriteString(strings.Join(selectColumns, ", "))
	sql.WriteString(" FROM ")
	sql.WriteString(GetFullTableName(queryTable))
	sql.WriteString(" WHERE ")
	sql.WriteString(quoteHistoryIdentifier(timestampField))
	sql.WriteString(" BETWEEN $1 AND $2")

	args := []any{time.UnixMilli(req.TimeStart), time.UnixMilli(req.TimeEnd)}
	if tbFieldName := strings.TrimSpace(uns.GetTbFieldName()); tbFieldName != "" && queryTable != uns.GetAlias() {
		sql.WriteString(" AND ")
		sql.WriteString(quoteHistoryIdentifier(tbFieldName))
		sql.WriteString(fmt.Sprintf(" = $%d", len(args)+1))
		args = append(args, uns.GetId())
	}

	sql.WriteString(" ORDER BY ")
	sql.WriteString(quoteHistoryIdentifier(timestampField))
	sql.WriteString(" ")
	sql.WriteString(normalizeHistoryOrder(req.Order))
	sql.WriteString(fmt.Sprintf(" LIMIT %d", normalizeHistoryLimit(req.Limit)))

	rows, err := dbPool.Query(ctx, sql.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fieldTypeMap := make(map[string]types.FieldType, len(uns.GetFields()))
	for _, field := range selectedFields {
		if field == nil {
			continue
		}
		fieldTypeMap[field.Name] = field.GetType()
	}

	result := &types.UnsHistoryFileResult{
		Alias:  uns.GetAlias(),
		Table:  queryTable,
		Fields: buildHistoryFields(selectedFields),
		List:   make([]*types.UnsHistoryValuePoint, 0),
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		row := make(map[string]any, len(values))
		for i, desc := range rows.FieldDescriptions() {
			row[desc.Name] = values[i]
		}

		timestamp, err := parseHistoryTimestamp(row[timestampField])
		if err != nil {
			return nil, err
		}

		payload := make(map[string]any, len(row)-1)
		for key, value := range row {
			if strings.EqualFold(key, timestampField) {
				continue
			}
			payload[key] = normalizeHistoryValue(value, fieldTypeMap[key])
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		result.List = append(result.List, &types.UnsHistoryValuePoint{
			Timestamp: timestamp.UnixMilli(),
			Payload:   string(payloadBytes),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func loadHistoryColumns(ctx context.Context, dbPool *pgxpool.Pool, queryTable string) (map[string]struct{}, error) {
	schemaName, tableName := splitHistoryTableName(queryTable)
	if tableName == "" {
		return nil, fmt.Errorf("history query table is empty")
	}

	sql := `SELECT column_name FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = $1`
	args := []any{tableName}
	if schemaName != "" {
		sql = `SELECT column_name FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2`
		args = []any{schemaName, tableName}
	}

	rows, err := dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make(map[string]struct{})
	for rows.Next() {
		var columnName string
		if err = rows.Scan(&columnName); err != nil {
			return nil, err
		}
		columns[columnName] = struct{}{}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return columns, nil
}

func splitHistoryTableName(queryTable string) (schemaName string, tableName string) {
	queryTable = strings.TrimSpace(queryTable)
	if queryTable == "" {
		return "", ""
	}
	parts := strings.Split(queryTable, ".")
	if len(parts) == 1 {
		return "", GetCleanTableName(parts[0])
	}
	return GetCleanTableName(strings.Join(parts[:len(parts)-1], ".")), GetCleanTableName(parts[len(parts)-1])
}

func buildHistorySelectColumns(uns types.UnsInfo, availableColumns map[string]struct{}) ([]string, []*types.FieldDefine) {
	columns := make([]string, 0, len(uns.GetFields()))
	selectedFields := make([]*types.FieldDefine, 0, len(uns.GetFields()))
	for _, field := range uns.GetFields() {
		if field == nil {
			continue
		}
		if len(availableColumns) > 0 {
			if _, ok := availableColumns[field.Name]; !ok {
				continue
			}
		}
		columns = append(columns, quoteHistoryIdentifier(field.Name))
		selectedFields = append(selectedFields, field)
	}
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	return columns, selectedFields
}

func buildHistoryFields(fieldDefines []*types.FieldDefine) []*types.UnsHistoryField {
	fields := make([]*types.UnsHistoryField, 0, len(fieldDefines))
	for _, field := range fieldDefines {
		if field == nil {
			continue
		}
		item := &types.UnsHistoryField{
			Name: field.Name,
			Type: field.Type,
		}
		if field.Unit != nil {
			item.Unit = *field.Unit
		}
		fields = append(fields, item)
	}
	return fields
}

func normalizeHistoryLimit(limit int) int {
	if limit <= 0 {
		return defaultHistoryQueryLimit
	}
	if limit > maxHistoryQueryLimit {
		return maxHistoryQueryLimit
	}
	return limit
}

func normalizeHistoryOrder(order string) string {
	if strings.EqualFold(order, "DESC") {
		return "DESC"
	}
	return "ASC"
}

func quoteHistoryIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func normalizeHistoryValue(value any, fieldType types.FieldType) any {
	switch v := value.(type) {
	case time.Time:
		if fieldType == types.FieldTypeDatetime {
			return v.UnixMilli()
		}
		return v.Format(time.RFC3339Nano)
	case *time.Time:
		if v == nil {
			return nil
		}
		return normalizeHistoryValue(*v, fieldType)
	case []byte:
		return string(v)
	default:
		return value
	}
}

func parseHistoryTimestamp(raw any) (time.Time, error) {
	switch v := raw.(type) {
	case time.Time:
		return v, nil
	case *time.Time:
		if v == nil {
			return time.Time{}, fmt.Errorf("history timestamp is nil")
		}
		return *v, nil
	case string:
		if parsed, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return parsed, nil
		}
		if parsed, err := time.Parse("2006-01-02 15:04:05.999999-07:00", v); err == nil {
			return parsed, nil
		}
		if parsed, err := time.Parse("2006-01-02 15:04:05-07:00", v); err == nil {
			return parsed, nil
		}
		return time.Time{}, fmt.Errorf("history timestamp parse failed: %s", v)
	case []byte:
		return parseHistoryTimestamp(string(v))
	case int64:
		return time.UnixMilli(v), nil
	case int32:
		return time.UnixMilli(int64(v)), nil
	case int:
		return time.UnixMilli(int64(v)), nil
	case float64:
		return time.UnixMilli(int64(v)), nil
	default:
		return time.Time{}, fmt.Errorf("history timestamp type invalid: %T", raw)
	}
}
