package event

import (
	"backend/internal/common/dto"
	"backend/internal/types"
)

// SaveDataEvent defines an event for saving data to a specified data source.
type SaveDataEvent struct {
	ApplicationEvent
	JdbcType        types.SrcJdbcType
	TopicData       []*dto.SaveDataDto
	DuplicateIgnore *bool
}
