package types

import "sync"

type UnsDefinition struct {
	CreateTopicDto
	Lock sync.RWMutex `json:"-"`
}
