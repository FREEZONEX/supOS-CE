package timescaledb

import "backend/internal/types"

// SQLGenerator SQL 生成器
type SQLGenerator struct{}

// NewSQLGenerator 创建 SQL 生成器
func NewSQLGenerator() *SQLGenerator {
	return &SQLGenerator{}
}

// RequiredFields 需要的字段信息
type RequiredFields struct {
	FieldsToAdd  []string
	FieldTypeMap map[string]types.FieldType
}
type ViewColumnInfo struct {
	ColumnName   string `json:"column_name"`   // 视图中的字段名
	SourceColumn string `json:"source_column"` // 物理表 uns_timeserial 中的字段名
}
type SimpleViewInfo struct {
	SrcTable string
	Columns  []ViewColumnInfo
}

// UnsViewInfo 包含 Uns 信息和视图信息
type UnsViewInfo struct {
	Uns  types.UnsInfo
	View SimpleViewInfo
}

// SyncSQLs 同步 SQL 结果
type SyncSQLs struct {
	CreateTableSQL []string            // 创建表的 SQL
	AlterTableSQL  []string            // 修改表的 SQL
	UpdateDataSQL  []string            // 更新数据的 SQL
	CreateViewSQL  []string            // 创建视图的 SQL
	FieldMappings  map[string][]string // 每个视图的字段映射
	Errors         []error             // 错误信息
}

// NewSyncSQLs 创建 SyncSQLs
func NewSyncSQLs() *SyncSQLs {
	return &SyncSQLs{
		CreateTableSQL: make([]string, 0),
		AlterTableSQL:  make([]string, 0),
		UpdateDataSQL:  make([]string, 0),
		CreateViewSQL:  make([]string, 0),
		FieldMappings:  make(map[string][]string),
		Errors:         make([]error, 0),
	}
}
