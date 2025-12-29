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

type MarkTopLogic struct {
	logx.Logger
	ctx                 context.Context
	svcCtx              *svc.ServiceContext
	dashboardMarkMapper relationDB.DashboardMarkedMapper
}

func NewMarkTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkTopLogic {
	return &MarkTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkTopLogic) MarkTop(id string, userID string) (*types.JsonResult, error) {
	mark := &relationDB.DashboardMarkModel{
		ID:     id,
		UserID: userID,
	}
	err := l.dashboardMarkMapper.Insert(l.ctx, mark)
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
