package common

import "backend/internal/common/constants"

// SrcJdbcType contains information about a JDBC type
type SrcJdbcType struct {
	ID          int    // ID
	DataSrcType string // Data source type
	Alias       string // Alias
	TypeCode    int    // Type code (1--时序，2--关系)
}

var (
	SrcJdbcTypeNone = SrcJdbcType{
		ID:          0,
		DataSrcType: "",
		Alias:       "",
		TypeCode:    0,
	}
	SrcJdbcTypeTdEngine = SrcJdbcType{
		ID:          1,
		DataSrcType: "tdengine-datasource",
		Alias:       "td",
		TypeCode:    constants.TimeSequenceType,
	}
	SrcJdbcTypePostgresql = SrcJdbcType{
		ID:          2,
		DataSrcType: "postgresql",
		Alias:       "pg",
		TypeCode:    constants.RelationType,
	}
	SrcJdbcTypeTimeScaleDB = SrcJdbcType{
		ID:          3,
		DataSrcType: "postgresql",
		Alias:       "tmsc",
		TypeCode:    constants.TimeSequenceType,
	}
)

var srcJdbcTypes = map[int]SrcJdbcType{
	SrcJdbcTypeNone.ID:        SrcJdbcTypeNone,
	SrcJdbcTypeTdEngine.ID:    SrcJdbcTypeTdEngine,
	SrcJdbcTypePostgresql.ID:  SrcJdbcTypePostgresql,
	SrcJdbcTypeTimeScaleDB.ID: SrcJdbcTypeTimeScaleDB,
}

// GetByID returns SrcJdbcType by ID
func GetSrcJdbcTypeByID(id int) SrcJdbcType {
	if srcJdbcType, ok := srcJdbcTypes[id]; ok {
		return srcJdbcType
	}
	return SrcJdbcTypeNone
}

// String returns the alias string representation
func (s SrcJdbcType) String() string {
	return s.Alias
}
