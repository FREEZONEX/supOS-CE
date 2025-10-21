package event

import (
	"backend/internal/common"
	"backend/internal/common/dto"
)

// InitTopicsEvent defines an event for initializing topics for different data sources.
type InitTopicsEvent struct {
	Topics map[common.SrcJdbcType][]*dto.CreateTopicDto
}
