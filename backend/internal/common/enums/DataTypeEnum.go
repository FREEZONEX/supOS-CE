package enums

const (
	TEMPLATE_TYPE         = "TEMPLATE_TYPE"         //(0, "模板"),
	TIME_SEQUENCE_TYPE    = "TIME_SEQUENCE_TYPE"    //(1, "时序"),
	RELATION_TYPE         = "RELATION_TYPE"         //(2, "关系"),
	CALCULATION_REAL_TYPE = "CALCULATION_REAL_TYPE" //(3,"实时计算"),
	CALCULATION_HIST_TYPE = "CALCULATION_HIST_TYPE" //(4, "历史值计算"),
	ALARM_RULE_TYPE       = "ALARM_RULE_TYPE"       //(5,"报警规则类型"),
	MERGE_TYPE            = "MERGE_TYPE"            //(6, "聚合类型"),
	CITING_TYPE           = "CITING_TYPE"           //(7,"引用类型，不持久化，只读, 不能引用引用类型的文件"),
	JSONB_TYPE            = "JSONB_TYPE"            //(8, "JSONB 整个json当做一个字段存储"),
)

var dataTypes = make([]string, 10)
var dataTypeNames = make(map[string]int16, 16)

func init() {
	dataTypes[0] = TEMPLATE_TYPE
	dataTypes[1] = TIME_SEQUENCE_TYPE
	dataTypes[2] = RELATION_TYPE
	dataTypes[3] = CALCULATION_REAL_TYPE
	dataTypes[4] = CALCULATION_HIST_TYPE
	dataTypes[5] = ALARM_RULE_TYPE
	dataTypes[6] = MERGE_TYPE
	dataTypes[7] = CITING_TYPE
	dataTypes[8] = JSONB_TYPE

	for i, dt := range dataTypes {
		if len(dt) > 0 {
			dataTypeNames[dt] = int16(i)
		}
	}
}
func DataTypeName(dataType int16) string {
	if dataType >= 0 && int(dataType) < len(dataTypes) {
		return dataTypes[dataType]
	}
	return ""
}
func DataTypeInt(name string) int16 {
	if v, has := dataTypeNames[name]; has {
		return v
	}
	return -1
}
