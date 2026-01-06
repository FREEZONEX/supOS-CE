package dashboard

import (
	"backend/internal/common/dto"
	"backend/internal/common/errors"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageListLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	dashboardMapper relationDB.DashboardMapper
}

func NewPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageListLogic {
	return &PageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageListLogic) PageList(req *types.PageListRequest, userID string) (*dto.PageResultDTO[*relationDB.DashboardExtends], error) {
	orderCode := req.OrderCode
	if req.PageNo < 1 {
		req.PageNo = 1
	}

	pageResult := &dto.PageResultDTO[*relationDB.DashboardExtends]{
		Code:     http.StatusOK,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}

	l.Logger.Debugf("PageListLogic: PageList request: %+v, userID: %s", req, userID)

	// 排序字段校验
	if orderCode != "" {
		if orderCode != "name" && orderCode != "createTime" {
			return nil, errors.NewBuzError(400, "illegal sort param")
		}
		// 驼峰转下划线
		orderCode = camelToSnake(orderCode)
	}
	db := relationDB.GetDb(l.ctx)
	// 查询总数
	var total int64
	// 查询数据
	dashboards, err := l.dashboardMapper.SelectDashboard(db, userID, req.K, req.Type, orderCode, req.IsAsc, req.PageNo, req.PageSize, &total)
	if err != nil {
		l.Logger.Error("查询Dashboard失败:", err)
		return pageResult, nil
	}
	pageResult.Total = total
	pageResult.Data = dashboards

	return pageResult, nil
}

// camelToSnake 驼峰转下划线
func camelToSnake(s string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	return strings.ToLower(re.ReplaceAllString(s, "${1}_${2}"))
}
