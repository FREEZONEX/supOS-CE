package dashboard

import (
	"backend/internal/logic/supos/uns/dashboard/dao"
	"backend/internal/logic/supos/uns/dashboard/model"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DashboardIDReq struct {
	ID string `path:"id"`
}

type GetByIdLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	dashboardMapper *dao.DashboardMapper
}

func NewGetByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByIdLogic {
	db := relationDB.GetDb(ctx)
	return &GetByIdLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		dashboardMapper: dao.NewDashboardMapper(db, ctx),
	}
}

func (l *GetByIdLogic) GetById(req *DashboardIDReq) (*model.DashboardModel, error) {
	return l.dashboardMapper.SelectById(req.ID)
}
