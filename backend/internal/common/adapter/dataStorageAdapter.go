package adapter

import (
	"backend/internal/types"
	"database/sql"
)

type DataStorageAdapter interface {
	Adapter

	GetJdbcType() types.SrcJdbcType
	GetJdbcTemplate() *sql.DB
	GetDataSourceProperties() DataSourceProperties
}
