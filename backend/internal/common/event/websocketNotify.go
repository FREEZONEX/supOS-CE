package event

// WebsocketNotifyEvent defines an event for websocket notifications.
type WebsocketNotifyEvent struct {
	UnsID int64
	Path  string
}
