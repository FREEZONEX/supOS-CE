package dashboard

import (
	"backend/internal/common/utils/grafanautil"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type EditLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	dashboardMapper relationDB.DashboardMapper
}

func NewEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditLogic {
	return &EditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditLogic) Edit(dashboard *relationDB.DashboardModel) (*types.JsonResult, error) {
	// 检查 Dashboard 是否存在
	db := relationDB.GetDb(l.ctx)
	existing, err := l.dashboardMapper.SelectById(db, dashboard.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &types.JsonResult{
				Code: http.StatusBadRequest,
				Msg:  i18ns.LocalizeMsg("uns.dashboard.not.exit"),
			}, nil
		}
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.not.exit"),
		}, nil
	}
	if existing == nil {
		return &types.JsonResult{
			Code: http.StatusBadRequest,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.not.exit"),
		}, nil
	}

	// Grafana Dashboard 更新
	if existing.Type == 1 {
		// 获取现有的 Dashboard
		dbJSON, err := grafanautil.GetDashboardByUUID(l.ctx, dashboard.ID)
		if err != nil || dbJSON == nil {
			l.Logger.Errorf("failed to get grafana dashboard: %v", err)
			goto UPDATE_DB
		}

		// 更新 title 和 description
		if dashboardObj, ok := dbJSON["dashboard"].(map[string]any); ok {
			dashboardObj["title"] = dashboard.Name
			dashboardObj["description"] = dashboard.Description
		}

		// 调用 Grafana API 更新
		jsonBytes, _ := json.Marshal(dbJSON)
		url := grafanautil.GetGrafanaURL() + "/api/dashboards/db"
		_, err = grafanautil.CreateDashboardByBody(l.ctx, dashboard.ID, "", string(jsonBytes))
		if err != nil {
			l.Logger.Errorf("failed to update grafana dashboard: %v", err)
			return &types.JsonResult{
				Code: http.StatusInternalServerError,
				Msg:  i18ns.LocalizeMsg("uns.dashboard.edit.failed"),
			}, nil
		}
		l.Logger.Infof("updated grafana dashboard: %s, url: %s", dashboard.ID, url)
	}

UPDATE_DB:
	// 更新数据库
	dashboard.UpdateTime = time.Now()
	err = l.dashboardMapper.UpdateById(db, dashboard)
	if err != nil {
		l.Logger.Errorf("failed to update dashboard: %v", err)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.edit.failed"),
		}, nil
	}
	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "success",
	}, nil
}
