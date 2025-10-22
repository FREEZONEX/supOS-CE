package event

import "backend/internal/common/dto"

// UpdateInstanceEvent defines an event for updating UNS instances, folders, or templates.
type UpdateInstanceEvent struct {
	ApplicationEvent
	Topics    []*dto.CreateTopicDto
	Folder    []*dto.CreateTopicDto
	Templates []*dto.CreateTopicDto
}
