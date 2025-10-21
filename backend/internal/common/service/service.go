package service

import (
	"backend/internal/common/dto"
	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/dto/resource"
)

type IUnsDefinitionService interface {
	GetDefinitionByAlias(alias string) *dto.CreateTopicDto
	GetDefinitionByPath(path string) *dto.CreateTopicDto
	GetDefinitionById(id int64) *dto.CreateTopicDto

	DeleteByIds(ids []int64) error
	InitSaveBatch(list []*dto.CreateTopicDto, ver int64) error
	InitDeleteBatch(list []*dto.CreateTopicDto, ver int64) error
	SaveBatch(list []*dto.CreateTopicDto, ver int64) error
	DeleteBatch(list []*dto.CreateTopicDto, ver int64) error
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
	CreateModelAndInstance(topicDtos []*dto.CreateTopicDto, fromImport bool) (map[string]string, error)
}
