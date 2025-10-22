package event

// UnsFirstDataSavedEvent defines an event triggered when a UNS topic saves data for the first time.
type UnsFirstDataSavedEvent struct {
	ApplicationEvent
	UnsID    int64
	UnsFlags int
}
