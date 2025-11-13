package service

import (
	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/dto/resource"
	"backend/internal/types"
)

type IUnsDefinitionService interface {
	GetDefinitionByAlias(alias string) *types.CreateTopicDto
	GetDefinitionByPath(path string) *types.CreateTopicDto
	GetDefinitionById(id int64) *types.CreateTopicDto

	DeleteByIds(ids []int64) error
	SaveBatch(list []*types.CreateTopicDto) error
	DeleteBatch(list []*types.CreateTopicDto) error
}

type IResourceService interface {
	SaveByExternal(dto *resource.SaveResource4ExternalDto) (int64, error)
	DeleteByCode(code string) (bool, error)
	DeleteBySource(source string) (bool, error)
}

type IRoleService interface {
	GetRoleListByUserId(userID string) ([]*authdto.RoleDto, error)
}

type IUnsManagerService interface {
	CreateModelAndInstance(topicDtos []*types.CreateTopicDto, fromImport bool) (map[string]string, error)
}
