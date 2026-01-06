package serviceApi

import "backend/internal/types"

type IUnsDefinitionService interface {
	GetDefinitionByAlias(alias string) *types.CreateTopicDto
	GetDefinitionByPath(path string) *types.CreateTopicDto
	GetDefinitionById(id int64) *types.CreateTopicDto

	DeleteByIds(ids []int64) error
	SaveBatch(list []*types.CreateTopicDto) error
	DeleteBatch(list []*types.CreateTopicDto) error
}
