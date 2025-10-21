package event

import (
	"backend/internal/common"
	"backend/internal/common/dto"
)

// BatchCreateTableEvent defines an event for batch creating database tables.
type BatchCreateTableEvent struct {
	FromImport bool
	FlowName   string
	Topics     map[common.SrcJdbcType][]*dto.CreateTopicDto
}

// SetFlowName sets the flow name for the event.
func (e *BatchCreateTableEvent) SetFlowName(flowName string) *BatchCreateTableEvent {
	if flowName != "" {
		e.FlowName = flowName
	}
	return e
}
