package event

import (
	"backend/internal/logic/supos/uns/uns/bo"
	"context"
	"time"
)

// RemoveTopicsEvent defines a generic event for removing topics, templates, or folders.
type RemoveTopicsEvent struct {
	ApplicationEvent
	DeleteTime    time.Time
	WithFlow      bool
	WithDashboard bool
	Topics        []bo.UnsInfo
	Templates     []bo.UnsInfo
	Folders       []bo.UnsInfo
}

// NewRemoveTopicsEvent creates a new RemoveTopicsEvent.
func NewRemoveTopicsEvent(ctx context.Context, deleteTime time.Time, withFlow, withDashboard bool, topics, templates, folders []bo.UnsInfo) *RemoveTopicsEvent {
	return &RemoveTopicsEvent{
		ApplicationEvent: ApplicationEvent{Context: ctx},
		DeleteTime:       deleteTime,
		WithFlow:         withFlow,
		WithDashboard:    withDashboard,
		Topics:           topics,
		Templates:        templates,
		Folders:          folders,
	}
}
