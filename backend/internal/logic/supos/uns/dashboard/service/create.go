package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/dto/grafana"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	"backend/internal/repo/relationDB"
	"backend/internal/types"
	"context"
	"net/http"
	"time"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/google/uuid"
)

func (s *DashboardService) Create(ctx context.Context, req *relationDB.DashboardModel, creator string) (*types.JsonResult, error) {
	// 检查名称是否重复
	db := relationDB.GetDb(ctx)
	dashboardMapper := relationDB.DashboardMapper{}
	dashboards, _ := dashboardMapper.SelectByNameAndType(db, req.Name, req.Type)
	if len(dashboards) > 0 {
		return &types.JsonResult{
			Code: 500,
			Msg:  I18nUtils.GetMessage("uns.dashboard.name.duplicate"),
		}, nil
	}

	// 生成 ID
	req.ID = uuid.New().String()
	req.Creator = creator
	req.CreateTime = time.Now()
	req.UpdateTime = time.Now()

	// Grafana Dashboard 创建
	if req.Type == 1 {
		// 构建 Dashboard JSON, 只构建 dashboard 内部的对象
		template := grafanautil.LoadTemplate(ctx, "templates/dashboard-blank.json")

		pgParams := grafana.GrafanaDashboardParam{
			UID:   req.ID,
			Title: req.Name,
		}
		params := make(map[string]interface{})
		er := UnsConverter.CopyPropertiesDefault(pgParams, &params)
		if er != nil || len(params) == 0 {
			params["uid"] = req.ID
			params["title"] = req.Name
			params["version"] = 0
		}
		dashboardJson := grafanautil.FormatTemplateMap(template, params)
		s.logger.Info(">>>>>>>>>>>>>>>dashboardJson :{}", dashboardJson)

		// 调用 Grafana API 创建 Dashboard
		// grafanautil.CreateDashboardByBody 会自动添加 "dashboard": {} 和 "overwrite": true 的外层包装
		url := grafanautil.GetGrafanaURL() + "/api/dashboards/db"
		_, err := grafanautil.CreateDashboardByBody(ctx, req.ID, "", dashboardJson)
		if err != nil {
			s.logger.Errorf("failed to create grafana dashboard: %v", err)
			return &types.JsonResult{
				Code: http.StatusInternalServerError,
				Msg:  i18ns.LocalizeMsg("uns.dashboard.create.failed", err.Error()),
			}, nil
		}
		s.logger.Infof("created grafana dashboard: %s, url: %s", req.ID, url)
	}

	// 保存到数据库
	err := dashboardMapper.Insert(db, req)
	if err != nil {
		s.logger.Errorf("failed to save dashboard: %v", err)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.create.failed", err.Error()),
		}, nil
	}

	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "success",
		Data: req,
	}, nil
}
