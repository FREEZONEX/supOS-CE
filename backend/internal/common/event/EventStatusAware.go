package event

type EventStatusAware interface {
	BeforeEvent(totalListeners int, i int, listenerName string)
	AfterEvent(totalListeners int, i int, listenerName string, err error)
}
