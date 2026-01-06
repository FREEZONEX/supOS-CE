package serviceApi

import "backend/internal/common/dto/resource"

type IResourceService interface {
	SaveByExternal(dto *resource.SaveResource4ExternalDto) (int64, error)
	DeleteByCode(code string) (bool, error)
	DeleteBySource(source string) (bool, error)
}
