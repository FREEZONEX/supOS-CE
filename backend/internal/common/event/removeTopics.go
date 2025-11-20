package event

import (
	"backend/internal/types"
	"context"
	"time"
)

// RemoveTopicsEvent defines a generic event for removing topics, templates, or folders.
type RemoveTopicsEvent struct {
	ApplicationEvent
	DeleteTime    time.Time
	WithFlow      bool
	WithDashboard bool
	Topics        []*types.CreateTopicDto
	Templates     []*types.CreateTopicDto
	Folders       []*types.CreateTopicDto
}

// NewRemoveTopicsEvent creates a new RemoveTopicsEvent.
func NewRemoveTopicsEvent(ctx context.Context, deleteTime time.Time, withFlow, withDashboard bool, topics, templates, folders []*types.CreateTopicDto) *RemoveTopicsEvent {
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
