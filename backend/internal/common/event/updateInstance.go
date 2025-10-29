package event

import (
	"backend/internal/types"
)

// UpdateInstanceEvent defines an event for updating UNS instances, folders, or templates.
type UpdateInstanceEvent struct {
	ApplicationEvent
	Topics    []*types.CreateTopicDto
	Folder    []*types.CreateTopicDto
	Templates []*types.CreateTopicDto
}
