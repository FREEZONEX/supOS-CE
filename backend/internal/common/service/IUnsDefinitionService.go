package service

import "backend/internal/common/dto"

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
