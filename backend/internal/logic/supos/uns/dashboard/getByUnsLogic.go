package dashboard

import (
	unsservice "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetByUnsLogic struct {
	logx.Logger
	ctx                context.Context
	svcCtx             *svc.ServiceContext
	dashboardRefMapper *relationDB.DashboardRefMapper
	unsQueryService    *unsservice.UnsQueryService
}

func NewGetByUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByUnsLogic {
	db := relationDB.GetDb(ctx)
	// Note: UnsQueryService might be managed by spring, adjust if needed
	unsQueryService := &unsservice.UnsQueryService{}

	return &GetByUnsLogic{
		Logger:             logx.WithContext(ctx),
		ctx:                ctx,
		svcCtx:             svcCtx,
		dashboardRefMapper: relationDB.NewDashboardRefMapper(db, ctx),
		unsQueryService:    unsQueryService,
	}
}

func (l *GetByUnsLogic) GetByUns(unsAlias string) (*relationDB.DashboardModel, error) {
	// TODO: As identified before, the DTO from GetModelDefinition lacks the 'Refers' field.
	// This logic is simplified until the UNS service provides the necessary details.
	return l.dashboardRefMapper.GetByUns(unsAlias)
}
