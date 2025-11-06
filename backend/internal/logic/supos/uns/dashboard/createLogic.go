package dashboard

import (
	"backend/internal/common/errors"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/logic/supos/uns/dashboard/dao"
	"backend/internal/logic/supos/uns/dashboard/model"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	dashboardMapper *dao.DashboardMapper
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	db := relationDB.GetDb(ctx)
	return &CreateLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		dashboardMapper: dao.NewDashboardMapper(db, ctx),
	}
}

func (l *CreateLogic) Create(req *model.DashboardModel, creator string) (*model.DashboardModel, error) {
	// 检查名称是否重复
	dashboards, err := l.dashboardMapper.SelectByFlowNames([]string{req.Name})
	if err != nil {
		return nil, err
	}
	if len(dashboards) > 0 {
		for _, db := range dashboards {
			if db.Type == req.Type {
				return nil, errors.NewBuzError(500, "uns.dashboard.name.duplicate")
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
		// 构建 Dashboard JSON
		dashboardJSON := fmt.Sprintf(`{
			"dashboard": {
				"uid": "%s",
				"title": "%s",
				"id": null
			}
		}`, req.ID, req.Name)

		// 调用 Grafana API 创建 Dashboard
		url := grafanautil.GetGrafanaURL() + "/api/dashboards/db"
		_, err := grafanautil.CreateDashboardByBody(req.ID, "", dashboardJSON)
		if err != nil {
			l.Logger.Errorf("failed to create grafana dashboard: %v", err)
			return nil, errors.NewBuzError(500, "uns.dashboard.create.failed")
		}
		l.Logger.Infof("created grafana dashboard: %s, url: %s", req.ID, url)
	}

	// 保存到数据库
	err = l.dashboardMapper.Insert(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}
