package event

import (
	"backend/internal/types"
	"context"
	"time"
)

// TopicMessageEvent defines an event carrying a message from a topic.
type TopicMessageEvent struct {
	ApplicationEvent
	Def          *types.CreateTopicDto
	UnsID        int64
	DataType     int
	Topic        string
	Payload      string
	Protocol     string
	Data         map[string]any
	LastData     map[string]any
	LastDataTime map[string]int64
	FieldsMap    map[string]*types.FieldDefine
	NowInMills   int64
	Err          string
}

// NewTopicMessageEventSimple creates a new TopicMessageEvent with basic info.
func NewTopicMessageEventSimple(ctx context.Context, def *types.CreateTopicDto, unsID int64, dataType int, topic, payload string) *TopicMessageEvent {
	return &TopicMessageEvent{
		ApplicationEvent: ApplicationEvent{Context: ctx},
		Def:              def,
		UnsID:            unsID,
		DataType:         dataType,
		Topic:            topic,
		Payload:          payload,
		NowInMills:       time.Now().UnixMilli(),
	}
}
