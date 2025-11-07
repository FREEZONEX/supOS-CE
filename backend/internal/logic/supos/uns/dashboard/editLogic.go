package dashboard

import (
	"backend/internal/common/errors"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type EditLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	dashboardMapper *relationDB.DashboardMapper
}

func NewEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditLogic {
	db := relationDB.GetDb(ctx)
	return &EditLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		dashboardMapper: relationDB.NewDashboardMapper(db, ctx),
	}
}

func (l *EditLogic) Edit(dashboard *relationDB.DashboardModel) error {
	// 检查 Dashboard 是否存在
	existing, err := l.dashboardMapper.SelectById(dashboard.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewBuzError(400, "uns.dashboard.not.exit")
		}
		return err
	}
	if existing == nil {
		return errors.NewBuzError(400, "uns.dashboard.not.exit")
	}

	// Grafana Dashboard 更新
	if existing.Type == 1 {
		// 获取现有的 Dashboard
		dbJSON, err := grafanautil.GetDashboardByUUID(dashboard.ID)
		if err != nil || dbJSON == nil {
			return errors.NewBuzError(400, "uns.dashboard.not.exit")
		}

		// 更新 title 和 description
		if dashboardObj, ok := dbJSON["dashboard"].(map[string]any); ok {
			dashboardObj["title"] = dashboard.Name
			dashboardObj["description"] = dashboard.Description
		}

		// 调用 Grafana API 更新
		jsonBytes, _ := json.Marshal(dbJSON)
		url := grafanautil.GetGrafanaURL() + "/api/dashboards/db"
		_, err = grafanautil.CreateDashboardByBody(dashboard.ID, "", string(jsonBytes))
		if err != nil {
			l.Logger.Errorf("failed to update grafana dashboard: %v", err)
			return err
		}
		l.Logger.Infof("updated grafana dashboard: %s, url: %s", dashboard.ID, url)
	}

	// 更新数据库
	dashboard.UpdateTime = time.Now()
	return l.dashboardMapper.UpdateById(dashboard)
}
