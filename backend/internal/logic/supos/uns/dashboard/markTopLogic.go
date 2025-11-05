// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"backend/internal/logic/supos/uns/dashboard/dao"
	"backend/internal/logic/supos/uns/dashboard/model"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkTopLogic struct {
	logx.Logger
	ctx                 context.Context
	svcCtx              *svc.ServiceContext
	dashboardMarkMapper *dao.DashboardMarkedMapper
}

func NewMarkTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkTopLogic {
	db := relationDB.GetDb(ctx)
	return &MarkTopLogic{
		Logger:              logx.WithContext(ctx),
		ctx:                 ctx,
		svcCtx:              svcCtx,
		dashboardMarkMapper: dao.NewDashboardMarkedMapper(db, ctx),
	}
}

func (l *MarkTopLogic) MarkTop(id string, userID string) error {
	mark := &model.DashboardMarkModel{
		ID:     id,
		UserID: userID,
	}
	return l.dashboardMarkMapper.Insert(mark)
}
