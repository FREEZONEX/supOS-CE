package event

// UnsTopologyChangeEvent defines a marker event for UNS topology changes.
// It does not contain any data.
type UnsTopologyChangeEvent struct {
	ApplicationEvent
	TopologyMsg []byte
}
