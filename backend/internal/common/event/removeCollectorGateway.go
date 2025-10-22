package event

// RemoveCollectorGatewayEvent defines an event for removing a collector gateway.
type RemoveCollectorGatewayEvent struct {
	ApplicationEvent
	AuthUUID string
}
