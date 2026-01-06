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
	dashboardMapper relationDB.DashboardMapper
}

func NewGetByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByIdLogic {
	return &GetByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetByIdLogic) GetById(req *DashboardIDReq) (*types.JsonResult, error) {
	db := relationDB.GetDb(l.ctx)
	dashboard, err := l.dashboardMapper.SelectById(db, req.ID)
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
