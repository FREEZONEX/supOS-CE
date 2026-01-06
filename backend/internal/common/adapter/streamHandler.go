package adapter

import "backend/internal/common/dto"

type StreamHandler interface {
	CreateStream(namedSQL map[string]string)
	DeleteStream(name string)
	ListByNames(names []string) ([]*dto.StreamInfo, error)
	StopStream(name string)
	ResumeStream(name string)
}
