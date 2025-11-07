package dashboard

import (
	"backend/internal/common/dto"
	"backend/internal/common/errors"
	"backend/internal/common/utils/dbutil"
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
	dashboardMapper *relationDB.DashboardMapper
}

func NewPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageListLogic {
	db := relationDB.GetDb(ctx)
	return &PageListLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		dashboardMapper: relationDB.NewDashboardMapper(db, ctx),
	}
}

func (l *PageListLogic) PageList(req *types.PageListRequest, userID string) (*dto.PageResultDTO[*relationDB.DashboardModel], error) {
	keyword := dbutil.EscapeForLike(req.K)
	orderCode := req.OrderCode
	descOrAsc := req.IsAsc

	pageResult := &dto.PageResultDTO[*relationDB.DashboardModel]{
		Code:     http.StatusOK,
		PageNo:   int64(req.PageNum),
		PageSize: int64(req.PageSize),
	}

	l.Logger.Infof("PageListLogic: PageList request: %+v, userID: %s", req, userID)

	// 排序字段校验
	if orderCode != "" {
		if orderCode != "name" && orderCode != "createTime" {
			return nil, errors.NewBuzError(400, "illegal sort param")
		}
		// 驼峰转下划线
		orderCode = camelToSnake(orderCode)
	}

	// 排序方向
	if descOrAsc == "" || descOrAsc != "ASC" {
		descOrAsc = "DESC"
	}

	// 查询总数
	total, err := l.dashboardMapper.SelectDashboardCount(keyword, req.Type)
	if err != nil {
		return nil, err
	}
	pageResult.Total = total

	// 查询数据
	dashboards, err := l.dashboardMapper.SelectDashboard(userID, keyword, req.Type, orderCode, descOrAsc, int64(req.PageNum), int64(req.PageSize))
	if err != nil {
		return nil, err
	}

	// 构建响应
	data := make([]*relationDB.DashboardModel, len(dashboards))
	for i, db := range dashboards {
		data[i] = &db.DashboardModel
	}
	pageResult.Data = data

	return pageResult, nil
}

// camelToSnake 驼峰转下划线
func camelToSnake(s string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	return strings.ToLower(re.ReplaceAllString(s, "${1}_${2}"))
}
