package event

import "backend/internal/common/enums"

// SysEvent defines a generic system event.
type SysEvent struct {
	Service   string
	EventMeta string
	Action    string
	Payload   any
}

// NewSysEventFromEnums creates a new SysEvent from enum types.
func NewSysEventFromEnums(service enums.ServiceEnum, eventMeta enums.EventMetaEnum, action enums.ActionEnum, payload any) *SysEvent {
	return &SysEvent{
		Service:   service.GetCode(),
		EventMeta: eventMeta.GetCode(),
		Action:    action.GetCode(),
		Payload:   payload,
	}
}
