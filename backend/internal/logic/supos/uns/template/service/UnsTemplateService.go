package service

import (
	"backend/internal/logic/supos/uns/uns/service"
	dao "backend/internal/repo/relationDB"
	"backend/share/spring"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsTemplateService struct {
	log              logx.Logger
	unsMapper        dao.UnsNamespaceRepo
	unsAddService    *service.UnsAddService
	unsRemoveService *service.UnsRemoveService
}

func init() {
	spring.RegisterLazy[*UnsTemplateService](func() *UnsTemplateService {
		return &UnsTemplateService{
			log:              logx.WithContext(context.Background()),
			unsAddService:    spring.GetBean[*service.UnsAddService](),
			unsRemoveService: spring.GetBean[*service.UnsRemoveService](),
		}
	})
}
