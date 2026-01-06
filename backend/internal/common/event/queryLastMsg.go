package event

import (
	"backend/internal/types"
)

// QueryLastMsgEvent defines an event for querying the last message of a topic.
// It holds the result in the LastMessage and MsgCreateTime fields.
type QueryLastMsgEvent struct {
	ApplicationEvent
	Uns           *types.CreateTopicDto
	MsgCreateTime int64
	LastMessage   map[string]any
}
