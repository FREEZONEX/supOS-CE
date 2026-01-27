// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupedDashboardListLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	dashboardMapper relationDB.DashboardMapper
}

// 分页按分组获取dashboard列表
func NewGetGroupedDashboardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupedDashboardListLogic {
	return &GetGroupedDashboardListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupedDashboardListLogic) GetGroupedDashboardList(req *types.GroupPageRequest) (resp *types.DashboardGroupPageResp, err error) {
	req.GroupType = 3
	if req.PageNo < 1 {
		req.PageNo = 1
	}

	pageResult := &types.DashboardGroupPageResp{
		PageResultDTO: types.PageResultDTO{
			Code:     http.StatusOK,
			PageNo:   req.PageNo,
			PageSize: req.PageSize,
		},
	}

	l.Logger.Debugf("GetGroupedDashboardList: PageList request: %+v", req)
	db := relationDB.GetDb(l.ctx)

	// 调用DAO层方法获取分组后的dashboard列表
	items, total, err := l.dashboardMapper.GetGroupedDashboardList(db, req)
	if err != nil {
		l.Logger.Error("查询分组dashboard列表失败:", err)
		return pageResult, nil
	}

	// 转换为前端需要的格式
	var groupBizVOList []types.GroupBizVO
	for _, item := range items {
		groupBizVO := types.GroupBizVO{}
		//存在分组
		if item.GroupType != 0 {
			// 处理分组数据
			groupID, _ := parseStringToInt64(item.ID)
			groupBizVO.ID = groupID
			groupBizVO.GroupType = &item.GroupType
			groupBizVO.Name = item.Name
			groupBizVO.Description = item.Description
			groupBizVO.Sort = item.Sort
			groupBizVO.CreateAt = item.CreateAt.UnixMilli()
			groupBizVO.GroupId = item.GroupID
			groupBizVO.Category = item.Category
		} else if item.GroupType == 0 {
			// 处理未分组的dashboard数据
			bizID, _ := parseStringToInt64(item.ID)
			groupBizVO.ID = bizID
			groupBizVO.Name = item.Name
			groupBizVO.Description = item.Description
			groupBizVO.CreateAt = item.CreateAt.UnixMilli()
			groupBizVO.GroupId = item.GroupID
			groupBizVO.Category = item.Category
		}
		groupBizVOList = append(groupBizVOList, groupBizVO)
	}

	pageResult.Total = total
	pageResult.Data = groupBizVOList

	return pageResult, nil
}

// parseStringToInt64 字符串转int64
func parseStringToInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
