package event

import (
	"backend/internal/common/enums"
	"context"
)

// SysEvent defines a generic system event.
type SysEvent struct {
	ApplicationEvent
	Service   string
	EventMeta string
	Action    string
	Payload   any
}

// NewSysEventFromEnums creates a new SysEvent from enum types.
func NewSysEventFromEnums(ctx context.Context, service enums.ServiceEnum, eventMeta enums.EventMetaEnum, action enums.ActionEnum, payload any) *SysEvent {
	return &SysEvent{
		ApplicationEvent: ApplicationEvent{Context: ctx},
		Service:          service.GetCode(),
		EventMeta:        eventMeta.GetCode(),
		Action:           action.GetCode(),
		Payload:          payload,
	}
}
