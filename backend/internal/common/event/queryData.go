package event

import "backend/internal/common/dto"

// QueryDataEvent defines an event for querying data.
// It holds the query conditions and is used to store the results.
type QueryDataEvent struct {
	ApplicationEvent
	TopicDto     *dto.CreateTopicDto
	EQConditions []*dto.EQCondition
	Values       []map[string]any
}
