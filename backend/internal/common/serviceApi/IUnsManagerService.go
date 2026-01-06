package serviceApi

import "backend/internal/types"

type IUnsManagerService interface {
	CreateModelAndInstance(topicDtos []*types.CreateTopicDto, fromImport bool) (map[string]string, error)
}
