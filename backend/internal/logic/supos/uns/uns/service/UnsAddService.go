package service

import (
	"backend/internal/logic/supos/uns/label/service"
	"backend/internal/logic/supos/uns/uns/bo"
	"backend/internal/repo/relationDB"
	"backend/share/spring"
	"context"
)

type UnsAddService struct {
	unsMapper       relationDB.UnsNamespaceRepo
	labelRefMapper  relationDB.UnsLabelRefRepo
	unsLabelService *service.UnsLabelService
}

func init() {
	spring.RegisterLazy[*UnsAddService](func() *UnsAddService {
		return &UnsAddService{
			unsLabelService: spring.GetBean[*service.UnsLabelService](),
		}
	})
}
func (u *UnsAddService) BatchCreateUns(ctx context.Context, args bo.CreateModelInstancesArgs) {

}
