package postgresql

import "backend/internal/types"

// 字段类型到数据库类型的映射
var fieldType2DBTypeMap = map[string]string{
	types.FieldTypeInteger:  "int4",
	types.FieldTypeLong:     "int8",
	types.FieldTypeFloat:    "float4",
	types.FieldTypeDouble:   "float8",
	types.FieldTypeString:   "text",
	types.FieldTypeBoolean:  "boolean",
	types.FieldTypeDatetime: "timestamptz",
	types.FieldTypeBlob:     "varchar(512)",
	types.FieldTypeLBlob:    "varchar(512)",
}

// 全局的类型映射表
var dbType2FieldTypeMap = make(map[string]string, 32)

// 初始化类型映射
func init() {
	// 整数类型
	dbType2FieldTypeMap["integer"] = types.FieldTypeInteger
	dbType2FieldTypeMap["serial"] = types.FieldTypeInteger
	dbType2FieldTypeMap["serial2"] = types.FieldTypeInteger
	dbType2FieldTypeMap["serial4"] = types.FieldTypeInteger
	dbType2FieldTypeMap["intserial"] = types.FieldTypeInteger
	dbType2FieldTypeMap["int2"] = types.FieldTypeInteger
	dbType2FieldTypeMap["int4"] = types.FieldTypeInteger

	// 长整型
	dbType2FieldTypeMap["bigint"] = types.FieldTypeLong
	dbType2FieldTypeMap["bigserial"] = types.FieldTypeLong
	dbType2FieldTypeMap["serial8"] = types.FieldTypeLong
	dbType2FieldTypeMap["int8"] = types.FieldTypeLong

	// 时间类型
	dbType2FieldTypeMap["timestamptz"] = types.FieldTypeDatetime
	dbType2FieldTypeMap["timestamp"] = types.FieldTypeDatetime

	// 浮点类型
	dbType2FieldTypeMap["float"] = types.FieldTypeFloat
	dbType2FieldTypeMap["float4"] = types.FieldTypeFloat
	dbType2FieldTypeMap["double"] = types.FieldTypeDouble
	dbType2FieldTypeMap["float8"] = types.FieldTypeDouble

	// 字符串类型
	dbType2FieldTypeMap["text"] = types.FieldTypeString
	dbType2FieldTypeMap["char"] = types.FieldTypeString
	dbType2FieldTypeMap["json"] = types.FieldTypeString
	dbType2FieldTypeMap["jsonb"] = types.FieldTypeString
	dbType2FieldTypeMap["varchar"] = types.FieldTypeString

	// 布尔类型
	dbType2FieldTypeMap["bool"] = types.FieldTypeBoolean
	dbType2FieldTypeMap["boolean"] = types.FieldTypeBoolean

	// 二进制类型
	dbType2FieldTypeMap["blob"] = types.FieldTypeBlob
}

// GetFieldType 根据数据库类型获取字段类型
func GetFieldType(dbType string) string {
	if fieldType, exists := dbType2FieldTypeMap[dbType]; exists {
		return fieldType
	}
	return types.FieldTypeString // 默认返回字符串类型
}
