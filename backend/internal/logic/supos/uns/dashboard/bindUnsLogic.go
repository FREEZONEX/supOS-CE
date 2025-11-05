package dashboard

import (
	"backend/internal/common/errors"
	"backend/internal/logic/supos/uns/dashboard/dao"
	"backend/internal/logic/supos/uns/dashboard/model"
	unsservice "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindUnsLogic struct {
	logx.Logger
	ctx                context.Context
	svcCtx             *svc.ServiceContext
	dashboardMapper    *dao.DashboardMapper
	dashboardRefMapper *dao.DashboardRefMapper
	unsQueryService    *unsservice.UnsQueryService
}

func NewBindUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUnsLogic {
	db := relationDB.GetDb(ctx)
	// Note: UnsQueryService might be managed by spring, adjust if needed
	unsQueryService := &unsservice.UnsQueryService{}

	return &BindUnsLogic{
		Logger:             logx.WithContext(ctx),
		ctx:                ctx,
		svcCtx:             svcCtx,
		dashboardMapper:    dao.NewDashboardMapper(db, ctx),
		dashboardRefMapper: dao.NewDashboardRefMapper(db, ctx),
		unsQueryService:    unsQueryService,
	}
}

func (l *BindUnsLogic) BindUns(dashboardID string, unsAlias string) error {
	// 检查 Dashboard 是否存在
	dashboard, err := l.dashboardMapper.SelectById(dashboardID)
	if err != nil {
		return err
	}
	if dashboard == nil {
		return errors.NewBuzError(400, "uns.dashboard.not.exit")
	}

	// 检查 UNS 是否存在
	unsResp, err := l.unsQueryService.GetModelDefinition(l.ctx, &types.ModelDetailReq{}, unsAlias)
	if err != nil {
		l.Logger.Errorf("failed to get uns definition for alias %s: %v", unsAlias, err)
		return errors.NewBuzError(500, "uns.file.not.exist")
	}
	if unsResp == nil || unsResp.Data == nil || unsResp.Data.Id == "" {
		return errors.NewBuzError(400, "uns.file.not.exist")
	}

	// 删除旧的绑定关系
	err = l.dashboardRefMapper.DeleteByDashboardId(dashboardID)
	if err != nil {
		l.Logger.Errorf("failed to delete old dashboard ref: %v", err)
	}

	// 创建新的绑定关系
	ref := &model.DashboardRefModel{
		DashboardID: dashboardID,
		UnsAlias:    unsAlias,
		CreateAt:    time.Now(),
	}
	return l.dashboardRefMapper.Insert(ref)
}
