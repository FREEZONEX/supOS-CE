package serviceApi

// TopicMessageConsumer 定义消息消费者接口
type TopicMessageConsumer interface {
	// OnMessageByAlias 处理单个别名消息
	OnMessageByAlias(alias, payload string)

	// OnBatchMessage 处理批量消息
	OnBatchMessage(payloads map[string]map[string]any)

	// OnMessageByAliasOnUpdate 处理vqt消息
	OnMessageByAliasOnUpdate(aliasVqtMap map[string]string)
}
