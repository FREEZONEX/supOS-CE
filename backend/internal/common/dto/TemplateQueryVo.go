package dto

// TemplateQueryVo 模板查询参数
type TemplateQueryVo struct {
	// Key 关键字查询，模版名称模糊匹配
	Key string `json:"key" schema:"关键字查询，模版名称模糊匹配"`

	// SubscribeEnable 是否订阅
	SubscribeEnable *bool `json:"subscribeEnable" schema:"是否订阅"`
}
