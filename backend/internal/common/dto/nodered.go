package dto

// CreateUnsNodeRedDto represents UNS Node-RED creation DTO
type CreateUnsNodeRedDto struct {
	Path      string `json:"path" validate:"required"` // 文件路径
	Alias     string `json:"alias,omitzero"`           // 别名
	FieldType string `json:"fieldType,omitzero"`       // 字段类型
	FieldName string `json:"fieldName,omitzero"`       // 字段名
	Tag       string `json:"tag,omitzero"`             // 标签
}

// NodeFlowDTO represents Node-RED flow DTO
type NodeFlowDTO struct {
	ID          string `json:"id,omitzero"`          // 流ID
	FlowName    string `json:"flowName,omitzero"`    // 流名称
	FlowID      string `json:"flowId,omitzero"`      // 流标识符
	Description string `json:"description,omitzero"` // 描述
	FlowStatus  string `json:"flowStatus,omitzero"`  // 流状态
	Template    string `json:"template,omitzero"`    // 模板
}

// NodeRedTagsDTO represents Node-RED tags DTO
type NodeRedTagsDTO struct {
	Data [][]string `json:"data,omitzero"` // 标签数据，每个元素是一个字符串数组
}
