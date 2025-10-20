package dto

// HistoryValueRequest represents history value query request
type HistoryValueRequest struct {
	Fill    *FillStrategy  `json:"fill,omitzero"`    // 补点策略
	GroupBy *GroupByWindow `json:"groupBy,omitzero"` // 聚合窗口
	Limit   int            `json:"limit,omitzero"`   // 返回的最大元素数目，默认1000，最大10000
	Offset  int            `json:"offset,omitzero"`  // 偏移量，从指定条数后开始查询
	Order   string         `json:"order,omitzero"`   // 排序方式：ASC（升序）、DESC（降序）
	Select  []string       `json:"select,omitzero"`  // 查询的字段表达式列表
	Where   map[string]any `json:"where,omitzero"`   // 过滤条件
}

// FillStrategy represents fill strategy for missing data points
type FillStrategy struct {
	Strategy string `json:"strategy,omitzero"` // 补点策略：none（不补点，默认）、previous（使用前一个窗口值补点）、line（线性补值）
}

// GroupByWindow represents aggregation window configuration
type GroupByWindow struct {
	Time string `json:"time,omitzero"` // 窗口配置字符串，格式：窗口间隔[,窗口偏移]，单位：s秒、m分、h小时、d天，例如 5s,1s
}

// NewHistoryValueRequest creates a new HistoryValueRequest with default values
func NewHistoryValueRequest() *HistoryValueRequest {
	return &HistoryValueRequest{
		Limit:  1000,
		Offset: 0,
		Order:  "desc",
	}
}

// FileBlobDataQueryDto represents file BLOB data query DTO
type FileBlobDataQueryDto struct {
	FileAlias    string         `json:"fileAlias,omitzero"`    // 文件别名
	EQConditions []*EQCondition `json:"eqConditions,omitzero"` // 相等条件列表
}

// EQCondition represents equality condition for BLOB query
type EQCondition struct {
	FieldName string `json:"fieldName,omitzero"` // 字段名
	Value     string `json:"value,omitzero"`     // 值
}

// NewFileBlobDataQueryDto creates a new FileBlobDataQueryDto
func NewFileBlobDataQueryDto(fileAlias string, conditions []*EQCondition) *FileBlobDataQueryDto {
	return &FileBlobDataQueryDto{
		FileAlias:    fileAlias,
		EQConditions: conditions,
	}
}

// NewEQCondition creates a new EQCondition
func NewEQCondition(fieldName, value string) *EQCondition {
	return &EQCondition{
		FieldName: fieldName,
		Value:     value,
	}
}
