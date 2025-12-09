package serviceApi

import "backend/internal/types"

type WebsocketMessage struct {
	Path string // 非 UNS 的消息，只有 Path 属性非空

	Def     *types.CreateTopicDto
	Data    map[string]any
	Payload string
	ErrMsg  string
}
type IWebsocketSender interface {
	SendMessage(msg WebsocketMessage)

	HasTopologies() bool
}
