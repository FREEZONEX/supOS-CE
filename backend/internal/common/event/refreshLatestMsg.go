package event

// RefreshLatestMsgEvent defines an event to refresh the latest message cache.
type RefreshLatestMsgEvent struct {
	ApplicationEvent
	UnsID    int64
	DataType int
	Path     string
	Payload  string
	Dt       map[string]int64
	Data     map[string]any
	ErrorMsg string
}
