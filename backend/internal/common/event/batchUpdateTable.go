package event

import (
	"backend/internal/common"
	"backend/internal/common/dto"
)

// BatchUpdateTableEvent defines an event for batch updating database tables.
type BatchUpdateTableEvent struct {
	Topics   []*dto.UpdateFieldDto
	JdbcType common.SrcJdbcType
}
