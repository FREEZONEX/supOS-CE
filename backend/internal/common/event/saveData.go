package event

import (
	"backend/internal/common"
	"backend/internal/common/dto"
)

// SaveDataEvent defines an event for saving data to a specified data source.
type SaveDataEvent struct {
	ApplicationEvent
	JdbcType        common.SrcJdbcType
	TopicData       []*dto.SaveDataDto
	DuplicateIgnore *bool
}
