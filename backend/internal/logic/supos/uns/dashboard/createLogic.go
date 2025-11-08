package dashboard

import (
	"backend/internal/common/utils/grafanautil"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"fmt"
	"net/http"
	"time"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	dashboardMapper *relationDB.DashboardMapper
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	db := relationDB.GetDb(ctx)
	return &CreateLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		dashboardMapper: relationDB.NewDashboardMapper(db, ctx),
	}
}

func (l *CreateLogic) Create(req *relationDB.DashboardModel, creator string) (*types.JsonResult, error) {
	// 检查名称是否重复
	dashboards, err := l.dashboardMapper.SelectByFlowNames([]string{req.Name})
	if err != nil {
		return nil, err
	}
	if len(dashboards) > 0 {
		for _, db := range dashboards {
			if db.Type == req.Type {
				return &types.JsonResult{
					Code: http.StatusBadRequest,
					Msg:  i18ns.LocalizeMsg("uns.dashboard.name.duplicate"),
				}, nil
			}
		}
	}

	// 生成 ID
	req.ID = uuid.New().String()
	req.Creator = creator
	req.CreateTime = time.Now()
	req.UpdateTime = time.Now()

	// Grafana Dashboard 创建
	if req.Type == 1 {
		// 构建 Dashboard JSON, 只构建 dashboard 内部的对象
		dashboardJSON := fmt.Sprintf(`{
			"uid": "%s",
			"title": "%s",
			"id": null
		}`, req.ID, req.Name)

		// 调用 Grafana API 创建 Dashboard
		// grafanautil.CreateDashboardByBody 会自动添加 "dashboard": {} 和 "overwrite": true 的外层包装
		url := grafanautil.GetGrafanaURL() + "/api/dashboards/db"
		_, err := grafanautil.CreateDashboardByBody(req.ID, "", dashboardJSON)
		if err != nil {
			l.Logger.Errorf("failed to create grafana dashboard: %v", err)
			return &types.JsonResult{
				Code: http.StatusInternalServerError,
				Msg:  i18ns.LocalizeMsg("uns.dashboard.create.failed", err.Error()),
			}, nil
		}
		l.Logger.Infof("created grafana dashboard: %s, url: %s", req.ID, url)
	}

	// 保存到数据库
	err = l.dashboardMapper.Insert(req)
	if err != nil {
		l.Logger.Errorf("failed to save dashboard: %v", err)
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
