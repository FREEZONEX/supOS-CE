package event

import (
	"backend/internal/types"
)

// BatchUpdateTableEvent defines an event for batch updating database tables.
type BatchUpdateTableEvent struct {
	ApplicationEvent
	Topics   []*types.UpdateFieldDto
	JdbcType types.SrcJdbcType
}
