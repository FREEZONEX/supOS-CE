package event

import (
	"backend/internal/types"
)

// InitTopicsEvent defines an event for initializing topics for different data sources.
type InitTopicsEvent struct {
	ApplicationEvent
	Topics map[types.SrcJdbcType][]*types.CreateTopicDto
}
