package dashboard

import (
	"backend/internal/common/dto"
	"backend/internal/common/errors"
	"backend/internal/common/utils/dbutil"
	"backend/internal/logic/supos/uns/dashboard/dao"
	"backend/internal/logic/supos/uns/dashboard/model"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"regexp"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageListLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	dashboardMapper *dao.DashboardMapper
}

func NewPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageListLogic {
	db := relationDB.GetDb(ctx)
	return &PageListLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		dashboardMapper: dao.NewDashboardMapper(db, ctx),
	}
}

func (l *PageListLogic) PageList(req *types.PageListRequest, userID string) (*dto.PageResultDTO[*model.DashboardModel], error) {
	keyword := dbutil.EscapeForLike(req.K)
	orderCode := req.OrderCode
	descOrAsc := req.IsAsc

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
	total, err := l.dashboardMapper.SelectDashboardCount(keyword, &req.Type)
	if err != nil {
		return nil, err
	}

	// 查询数据
	dashboards, err := l.dashboardMapper.SelectDashboard(userID, keyword, &req.Type, orderCode, descOrAsc, int64(req.PageNum), int64(req.PageSize))
	if err != nil {
		return nil, err
	}

	// 构建响应
	data := make([]*model.DashboardModel, len(dashboards))
	for i, db := range dashboards {
		data[i] = &db.DashboardModel
	}

	return &dto.PageResultDTO[*model.DashboardModel]{
		Code:     0,
		Total:    total,
		PageNo:   int64(req.PageNum),
		PageSize: int64(req.PageSize),
		Data:     data,
	}, nil
}

// camelToSnake 驼峰转下划线
func camelToSnake(s string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	return strings.ToLower(re.ReplaceAllString(s, "${1}_${2}"))
}
