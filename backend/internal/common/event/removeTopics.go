package event

import (
	"backend/internal/common/dto"
	"time"
)

// RemoveTopicsEvent defines a generic event for removing topics, templates, or folders.
type RemoveTopicsEvent struct {
	DeleteTime    time.Time
	WithFlow      bool
	WithDashboard bool
	Topics        []*dto.CreateTopicDto
	Templates     []*dto.CreateTopicDto
	Folders       []*dto.CreateTopicDto
}

// NewRemoveTopicsEvent creates a new RemoveTopicsEvent.
func NewRemoveTopicsEvent(deleteTime time.Time, withFlow, withDashboard bool, topics, templates, folders []*dto.CreateTopicDto) *RemoveTopicsEvent {
	if topics == nil {
		topics = make([]*dto.CreateTopicDto, 0)
	}
	if templates == nil {
		templates = make([]*dto.CreateTopicDto, 0)
	}
	if folders == nil {
		folders = make([]*dto.CreateTopicDto, 0)
	}
	return &RemoveTopicsEvent{
		DeleteTime:    deleteTime,
		WithFlow:      withFlow,
		WithDashboard: withDashboard,
		Topics:        topics,
		Templates:     templates,
		Folders:       folders,
	}
}
