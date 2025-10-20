package dto

// DropTableDto represents drop table DTO
type DropTableDto struct {
	Topic   string `json:"topic" validate:"required"` // 主题
	SrcType string `json:"srcType,omitzero"`          // 数据源类型：pg--PostgreSQL, td--TdEngine
}
