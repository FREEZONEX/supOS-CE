package event

// RemoveTdengineEvent defines an event for removing TDengine topics.
type RemoveTdengineEvent struct {
	ApplicationEvent
	*RemoveTopicsEvent
}
