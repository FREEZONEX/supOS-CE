package event

// NamespaceChangeEvent defines an event for namespace changes.
type NamespaceChangeEvent struct {
	ApplicationEvent
	Topic string
	Data  map[string]any
}
