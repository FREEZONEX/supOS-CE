package dashboard

import (
	unsservice "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"net/http"
	"time"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/zeromicro/go-zero/core/logx"
)

type BindUnsLogic struct {
	logx.Logger
	ctx                context.Context
	svcCtx             *svc.ServiceContext
	dashboardMapper    relationDB.DashboardMapper
	dashboardRefMapper relationDB.DashboardRefMapper
	unsQueryService    *unsservice.UnsQueryService
}

func NewBindUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUnsLogic {
	// Note: UnsQueryService might be managed by spring, adjust if needed
	unsQueryService := &unsservice.UnsQueryService{}

	return &BindUnsLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		unsQueryService: unsQueryService,
	}
}

func (l *BindUnsLogic) BindUns(dashboardID string, unsAlias string) (*types.JsonResult, error) {
	// 检查 Dashboard 是否存在
	db := relationDB.GetDb(l.ctx)
	dashboard, err := l.dashboardMapper.SelectById(db, dashboardID)
	if err != nil {
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}
	if dashboard == nil {
		return &types.JsonResult{
			Code: http.StatusBadRequest,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.not.exit"),
		}, nil
	}

	// 检查 UNS 是否存在
	unsResp, err := l.unsQueryService.GetModelDefinition(l.ctx, &types.ModelDetailReq{}, unsAlias)
	if err != nil {
		l.Logger.Errorf("failed to get uns definition for alias %s: %v", unsAlias, err)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}
	if unsResp == nil || unsResp.Data == nil || unsResp.Data.Id == "" {
		l.Logger.Errorf("uns definition for alias %s not found", unsAlias)
		return &types.JsonResult{
			Code: http.StatusBadRequest,
			Msg:  i18ns.LocalizeMsg("uns.file.not.exist"),
		}, nil
	}

	// 删除旧的绑定关系
	err = l.dashboardRefMapper.DeleteByUnsAlias(db, unsAlias)
	if err != nil {
		l.Logger.Errorf("failed to delete old dashboard ref: %v", err)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}

	// 创建新的绑定关系
	ref := &relationDB.DashboardRefModel{
		DashboardID: dashboardID,
		UnsAlias:    unsAlias,
		CreateAt:    time.Now(),
	}
	err = l.dashboardRefMapper.Insert(db, ref)
	if err != nil {
		l.Logger.Errorf("failed to create new dashboard ref: %v", err)
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
