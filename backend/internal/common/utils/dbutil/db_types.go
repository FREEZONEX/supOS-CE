package dbutil

import (
	"backend/internal/types"
)

// PostgreSQLTypeMap maps PostgreSQL DB types to our internal FieldType.
var PostgreSQLTypeMap = map[string]types.FieldType{
	"integer":     types.FieldTypeInteger,
	"serial":      types.FieldTypeInteger,
	"serial2":     types.FieldTypeInteger,
	"serial4":     types.FieldTypeInteger,
	"intserial":   types.FieldTypeInteger,
	"int2":        types.FieldTypeInteger,
	"int4":        types.FieldTypeInteger,
	"bigint":      types.FieldTypeLong,
	"bigserial":   types.FieldTypeLong,
	"serial8":     types.FieldTypeLong,
	"int8":        types.FieldTypeLong,
	"timestamptz": types.FieldTypeDatetime,
	"timestamp":   types.FieldTypeDatetime,
	"float":       types.FieldTypeFloat,
	"float4":      types.FieldTypeFloat,
	"double":      types.FieldTypeDouble,
	"float8":      types.FieldTypeDouble,
	"text":        types.FieldTypeString,
	"char":        types.FieldTypeString,
	"json":        types.FieldTypeString,
	"jsonb":       types.FieldTypeString,
	"varchar":     types.FieldTypeString,
	"bool":        types.FieldTypeBoolean,
	"boolean":     types.FieldTypeBoolean,
	"blob":        types.FieldTypeBlob,
}

// MariaDBTypeMap maps MariaDB DB types to our internal FieldType.
var MariaDBTypeMap = map[string]types.FieldType{
	"int":        types.FieldTypeInteger,
	"smallint":   types.FieldTypeInteger,
	"tinyint":    types.FieldTypeInteger,
	"mediumint":  types.FieldTypeLong,
	"bigint":     types.FieldTypeLong,
	"float":      types.FieldTypeFloat,
	"double":     types.FieldTypeDouble,
	"decimal":    types.FieldTypeDouble,
	"char":       types.FieldTypeString,
	"varchar":    types.FieldTypeString,
	"text":       types.FieldTypeString,
	"tinytext":   types.FieldTypeString,
	"mediumtext": types.FieldTypeString,
	"json":       types.FieldTypeString,
	"bit":        types.FieldTypeLong,
	"date":       types.FieldTypeDatetime,
	"time":       types.FieldTypeDatetime,
	"datetime":   types.FieldTypeDatetime,
	"timestamp":  types.FieldTypeDatetime,
	"tinyblob":   types.FieldTypeBlob,
	"blob":       types.FieldTypeBlob,
	"mediumblob": types.FieldTypeBlob,
	"longblob":   types.FieldTypeBlob,
}

// SQLServerTypeMap maps SQL Server DB types to our internal FieldType.
var SQLServerTypeMap = map[string]types.FieldType{
	"int":            types.FieldTypeInteger,
	"smallint":       types.FieldTypeInteger,
	"tinyint":        types.FieldTypeInteger,
	"bigint":         types.FieldTypeLong,
	"float":          types.FieldTypeFloat,
	"real":           types.FieldTypeFloat,
	"decimal":        types.FieldTypeDouble,
	"numeric":        types.FieldTypeDouble,
	"char":           types.FieldTypeString,
	"varchar":        types.FieldTypeString,
	"text":           types.FieldTypeString,
	"nchar":          types.FieldTypeString,
	"nvarchar":       types.FieldTypeString,
	"ntext":          types.FieldTypeString,
	"bit":            types.FieldTypeBoolean,
	"date":           types.FieldTypeDatetime,
	"time":           types.FieldTypeDatetime,
	"datetime":       types.FieldTypeDatetime,
	"datetime2":      types.FieldTypeDatetime,
	"datetimeoffset": types.FieldTypeDatetime,
	"binary":         types.FieldTypeBlob,
	"varbinary":      types.FieldTypeBlob,
}
