package event

// AlertEvent defines an alert event, which is a specialized TopicMessageEvent.
type AlertEvent struct {
	*TopicMessageEvent
}
