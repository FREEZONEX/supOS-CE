package dashboard

import (
	"backend/internal/common/utils/fuxautil"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"errors"
	"fmt"
	"net/http"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteLogic struct {
	logx.Logger
	ctx                 context.Context
	svcCtx              *svc.ServiceContext
	dashboardMapper     relationDB.DashboardMapper
	dashboardRefMapper  relationDB.DashboardRefMapper
	dashboardMarkMapper relationDB.DashboardMarkedMapper
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(uid string) (*types.JsonResult, error) {
	// 检查 Dashboard 是否存在
	db := relationDB.GetDb(l.ctx)
	dashboard, err := l.dashboardMapper.SelectById(db, uid)
	if err != nil && err != gorm.ErrRecordNotFound {
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.delete.failed", err.Error()),
		}, nil
	}
	if dashboard == nil {
		return &types.JsonResult{
			Code: http.StatusOK,
			Msg:  "success",
		}, nil
	}

	// Grafana Dashboard 删除
	if dashboard.Type == 1 {
		err := grafanautil.DeleteDashboard(l.ctx, uid)
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
		_, err := resty.New().R().Delete(url)
		if err != nil {
			l.Logger.Errorf("failed to delete fuxa dashboard: %v", err)
		}
	}

	// 删除置顶标记
	err = l.dashboardMarkMapper.DeleteById(l.ctx, uid)
	if err != nil {
		l.Logger.Errorf("failed to delete dashboard mark: %v", err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.JsonResult{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}, nil
		}
	}

	// 删除引用关系
	err = l.dashboardRefMapper.DeleteByDashboardId(db, uid)
	if err != nil {
		l.Logger.Errorf("failed to delete dashboard ref: %v", err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.JsonResult{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}, nil
		}
	}

	// 删除 Dashboard
	err = l.dashboardMapper.DeleteById(db, uid)
	if err != nil {
		l.Logger.Errorf("failed to delete dashboard: %v", err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.JsonResult{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}, nil
		}
	}
	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "success",
	}, nil
}
