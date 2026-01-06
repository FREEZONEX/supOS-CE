// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnmarkTopLogic struct {
	logx.Logger
	ctx                 context.Context
	svcCtx              *svc.ServiceContext
	dashboardMarkMapper relationDB.DashboardMarkedMapper
}

func NewUnmarkTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnmarkTopLogic {
	return &UnmarkTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnmarkTopLogic) UnmarkTop(id string, userID string) (*types.JsonResult, error) {
	err := l.dashboardMarkMapper.Delete(l.ctx, id, userID)
	if err != nil {
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}
	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "success",
	}, nil
}
