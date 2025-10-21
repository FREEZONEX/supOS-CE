package adapter

import (
	"backend/internal/common"
	"database/sql"
)

type DataStorageAdapter interface {
	Adapter

	GetJdbcType() common.SrcJdbcType
	GetJdbcTemplate() *sql.DB
	GetDataSourceProperties() DataSourceProperties
}
