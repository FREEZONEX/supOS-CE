// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnmarkTopLogic struct {
	logx.Logger
	ctx                 context.Context
	svcCtx              *svc.ServiceContext
	dashboardMarkMapper *relationDB.DashboardMarkedMapper
}

func NewUnmarkTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnmarkTopLogic {
	db := relationDB.GetDb(ctx)
	return &UnmarkTopLogic{
		Logger:              logx.WithContext(ctx),
		ctx:                 ctx,
		svcCtx:              svcCtx,
		dashboardMarkMapper: relationDB.NewDashboardMarkedMapper(db, ctx),
	}
}

func (l *UnmarkTopLogic) UnmarkTop(id string, userID string) error {
	return l.dashboardMarkMapper.Delete(id, userID)
}
