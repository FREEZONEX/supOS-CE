package sdk

import "backend/internal/common/event"

type WebsocketSender interface {
	SendLatestMsg(event *event.WebsocketNotifyEvent)
}

type UnsQueryApi interface {
	RefreshLatestMsg(event *event.RefreshLatestMsgEvent)
}
