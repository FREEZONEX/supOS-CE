package dashboard

import (
	"backend/internal/logic/supos/uns/dashboard/service"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *relationDB.DashboardModel, creator string) (*types.JsonResult, error) {
	return spring.GetBean[*service.DashboardService]().Create(l.ctx, req, creator)
}
