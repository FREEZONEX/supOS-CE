package serviceApi

import "backend/internal/types"

type IUnsDefinitionService interface {
	GetDefinitionByAlias(alias string) *types.UnsDefinition
	GetDefinitionByPath(path string) *types.UnsDefinition
	GetDefinitionById(id int64) *types.UnsDefinition

	DeleteByIds(ids []int64) error
	SaveBatch(list []*types.UnsDefinition) error
	DeleteBatch(list []*types.UnsDefinition) error
}
