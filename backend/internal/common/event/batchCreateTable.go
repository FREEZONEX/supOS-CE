package event

import (
	"backend/internal/common/constants"
	"backend/internal/types"
	"backend/share/base"
)

// BatchCreateTableEvent defines an event for batch creating database tables.
type BatchCreateTableEvent struct {
	ApplicationEvent
	FromImport    bool
	FlowName      string
	Creates       map[int16][]*types.CreateTopicDto
	Updates       map[int16][]*types.CreateTopicDto
	DelegateAware EventStatusAware
	createFiles   map[types.SrcJdbcType][]*types.CreateTopicDto
	updateFiles   map[types.SrcJdbcType][]*types.CreateTopicDto
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
func (e *BatchCreateTableEvent) GetCreateFiles(sinkDataSrcType types.SrcJdbcType) []*types.CreateTopicDto {
	if e.createFiles == nil && len(e.Creates) > 0 {
		files := e.Creates[constants.PathTypeFile]
		e.createFiles = base.GroupBy(files, pathTypeGroupSrcType)
	}
	if e.createFiles != nil {
		return e.createFiles[sinkDataSrcType]
	}
	return nil
}
func (e *BatchCreateTableEvent) GetAllCreateFiles() map[types.SrcJdbcType][]*types.CreateTopicDto {
	if e.createFiles == nil && len(e.Creates) > 0 {
		files := e.Creates[constants.PathTypeFile]
		e.createFiles = base.GroupBy(files, pathTypeGroupSrcType)
	}
	return e.createFiles
}
func (e *BatchCreateTableEvent) GetUpdateFiles(sinkDataSrcType types.SrcJdbcType) []*types.CreateTopicDto {
	if e.updateFiles == nil && len(e.Updates) > 0 {
		files := e.Updates[constants.PathTypeFile]
		e.updateFiles = base.GroupBy(files, pathTypeGroupSrcType)
	}
	if e.updateFiles != nil {
		return e.updateFiles[sinkDataSrcType]
	}
	return nil
}
func pathTypeGroupSrcType(e *types.CreateTopicDto) types.SrcJdbcType {
	return types.SrcJdbcType(e.DataSrcID)
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
