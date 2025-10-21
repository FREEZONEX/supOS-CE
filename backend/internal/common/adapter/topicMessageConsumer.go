package adapter

type TopicMessageConsumer interface {
	OnMessageByAlias(alias string, payload string)
	OnBatchMessage(payloads map[string]map[string]any)
	OnMessageByAliasOnUpdate(aliasVqtMap map[string]string)
}
