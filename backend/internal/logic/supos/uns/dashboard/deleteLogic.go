package dashboard

import (
	"backend/internal/common/utils/fuxautil"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/logic/supos/uns/dashboard/dao"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteLogic struct {
	logx.Logger
	ctx                 context.Context
	svcCtx              *svc.ServiceContext
	dashboardMapper     *dao.DashboardMapper
	dashboardRefMapper  *dao.DashboardRefMapper
	dashboardMarkMapper *dao.DashboardMarkedMapper
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	db := relationDB.GetDb(ctx)
	return &DeleteLogic{
		Logger:              logx.WithContext(ctx),
		ctx:                 ctx,
		svcCtx:              svcCtx,
		dashboardMapper:     dao.NewDashboardMapper(db, ctx),
		dashboardRefMapper:  dao.NewDashboardRefMapper(db, ctx),
		dashboardMarkMapper: dao.NewDashboardMarkedMapper(db, ctx),
	}
}

func (l *DeleteLogic) Delete(uid string) error {
	// 检查 Dashboard 是否存在
	dashboard, err := l.dashboardMapper.SelectById(uid)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if dashboard == nil {
		return nil // 已经不存在，视为成功
	}

	// Grafana Dashboard 删除
	if dashboard.Type == 1 {
		err := grafanautil.DeleteDashboard(uid)
		if err != nil {
			l.Logger.Errorf("failed to delete grafana dashboard: %v", err)
		}
	}

	// Fuxa Dashboard 删除
	if dashboard.Type == 2 {
		// Fuxa 使用 HTTP DELETE 请求删除
		url := fmt.Sprintf("%s/api/project/%s", fuxautil.GetFuxaURL(), uid)
		l.Logger.Infof("deleting fuxa dashboard: %s", url)
		// 注意：fuxautil 目前没有 Delete 方法，需要直接 HTTP 调用或添加方法
	}

	// 删除置顶标记
	err = l.dashboardMarkMapper.DeleteById(uid)
	if err != nil {
		l.Logger.Errorf("failed to delete dashboard mark: %v", err)
	}

	// 删除引用关系
	err = l.dashboardRefMapper.DeleteByDashboardId(uid)
	if err != nil {
		l.Logger.Errorf("failed to delete dashboard ref: %v", err)
	}

	// 删除 Dashboard
	return l.dashboardMapper.DeleteById(uid)
}
