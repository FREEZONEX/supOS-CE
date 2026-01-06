package event

import (
	"backend/internal/common/dto"
	"backend/internal/types"
)

// QueryDataEvent defines an event for querying data.
// It holds the query conditions and is used to store the results.
type QueryDataEvent struct {
	ApplicationEvent
	TopicDto     *types.CreateTopicDto
	EQConditions []*dto.EQCondition
	Values       []map[string]any
}
