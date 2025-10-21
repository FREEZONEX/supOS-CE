package common

import "backend/internal/common/constants"

type SrcJdbcType int
type srcJdbcTypeInfo struct {
	dataSrcType string // Data source type
	alias       string // alias
	typeCode    int    // Type code (1--时序，2--关系)
}

const (
	SrcJdbcTypeNone SrcJdbcType = iota
	SrcJdbcTypeTdEngine
	SrcJdbcTypePostgresql
	SrcJdbcTypeTimeScaleDB
	SrcJdbcTypeVxBase
)

var srcJdbcTypes = map[SrcJdbcType]srcJdbcTypeInfo{
	SrcJdbcTypeNone: {
		dataSrcType: "",
		alias:       "",
		typeCode:    0,
	},
	SrcJdbcTypeTdEngine: {
		dataSrcType: "tdengine-datasource",
		alias:       "td",
		typeCode:    constants.TimeSequenceType,
	},
	SrcJdbcTypePostgresql: {
		dataSrcType: "postgresql",
		alias:       "pg",
		typeCode:    constants.RelationType,
	},
	SrcJdbcTypeTimeScaleDB: {
		dataSrcType: "postgresql",
		alias:       "tmsc",
		typeCode:    constants.TimeSequenceType,
	},
}

// GetByID returns SrcJdbcType by ID
func GetSrcJdbcTypeByID(id int) SrcJdbcType {
	k := SrcJdbcType(id)
	if _, ok := srcJdbcTypes[k]; ok {
		return k
	}
	return SrcJdbcTypeNone
}
func (s SrcJdbcType) Id() int {
	return int(s)
}
func (s SrcJdbcType) DataSrcType() string {
	return srcJdbcTypes[s].dataSrcType
}
func (s SrcJdbcType) Alias() string {
	return srcJdbcTypes[s].alias
}
func (s SrcJdbcType) TypeCode() int {
	return srcJdbcTypes[s].typeCode
}

// String returns the alias string representation
func (s SrcJdbcType) String() string {
	return srcJdbcTypes[s].alias
}
