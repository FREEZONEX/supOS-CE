package dashboard

import (
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type DashboardIDReq struct {
	ID string `path:"id"`
}

type GetByIdLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	dashboardMapper *relationDB.DashboardMapper
}

func NewGetByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByIdLogic {
	db := relationDB.GetDb(ctx)
	return &GetByIdLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		dashboardMapper: relationDB.NewDashboardMapper(db, ctx),
	}
}

func (l *GetByIdLogic) GetById(req *DashboardIDReq) (*types.JsonResult, error) {
	dashboard, err := l.dashboardMapper.SelectById(req.ID)
	if err != nil {
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}
	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "success",
		Data: dashboard,
	}, nil
}
