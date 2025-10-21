package event

import "backend/internal/common/dto"

// RemoveTimescaleTopicsEvent defines an event for removing TimescaleDB topics.
type RemoveTimescaleTopicsEvent struct {
	Standard    []dto.SimpleUnsInfo
	NonStandard []dto.SimpleUnsInfo
}
