package event

// NamespaceChangeEvent defines an event for namespace changes.
type NamespaceChangeEvent struct {
	Topic string
	Data  map[string]any
}
