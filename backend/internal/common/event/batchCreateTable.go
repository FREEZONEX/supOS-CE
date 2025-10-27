package event

import (
	"backend/internal/common"
	"backend/internal/common/dto"
)

// BatchCreateTableEvent defines an event for batch creating database tables.
type BatchCreateTableEvent struct {
	ApplicationEvent
	FromImport    bool
	FlowName      string
	Topics        map[common.SrcJdbcType][]*dto.CreateTopicDto
	Updates       []*dto.CreateTopicDto
	Folders       []*dto.CreateTopicDto
	Templates     []*dto.CreateTopicDto
	Labels        []*dto.CreateTopicDto
	DelegateAware EventStatusAware
}

// SetFlowName sets the flow name for the event.
func (e *BatchCreateTableEvent) SetFlowName(flowName string) *BatchCreateTableEvent {
	if flowName != "" {
		e.FlowName = flowName
	}
	return e
}
func (e *BatchCreateTableEvent) SetDelegateAware(delegateAware EventStatusAware) {
	e.DelegateAware = delegateAware
}
func (e *BatchCreateTableEvent) BeforeEvent(totalListeners int, i int, listenerName string) {
	if target := e.DelegateAware; target != nil {
		target.BeforeEvent(totalListeners, i, listenerName)
	}
}
func (e *BatchCreateTableEvent) AfterEvent(totalListeners int, i int, listenerName string, err error) {
	if target := e.DelegateAware; target != nil {
		target.AfterEvent(totalListeners, i, listenerName, err)
	}
}
