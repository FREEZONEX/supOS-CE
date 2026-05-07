// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"backend/internal/repo/iam"
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

	// 批量查询 creator 显示名
	creatorSet := make(map[string]struct{})
	for _, item := range items {
		if item.Creator != "" {
			creatorSet[item.Creator] = struct{}{}
		}
	}
	usernames := make([]string, 0, len(creatorSet))
	for u := range creatorSet {
		usernames = append(usernames, u)
	}
	displayNameMap := make(map[string]string)
	if len(usernames) > 0 {
		if authRepo, aErr := iam.NewAuthRepo(l.ctx); aErr == nil {
			displayNameMap, _ = authRepo.GetDisplayNamesByUsernames(usernames)
		}
	}

	for _, item := range items {
		groupBizVO := types.GroupBizVO{}
		groupBizVO.ID = item.ID
		groupBizVO.Name = item.Name
		groupBizVO.Description = item.Description
		groupBizVO.GroupType = &item.GroupType
		groupBizVO.Sort = item.Sort
		groupBizVO.CreateAt = item.CreateAt.UnixMilli()
		groupBizVO.Category = item.Category
		groupBizVO.Creator = item.Creator
		if dn := displayNameMap[item.Creator]; dn != "" {
			groupBizVO.Creator = dn
		}
		groupBizVO.HasChildren = item.HasChildren
		if item.GroupID != nil {
			groupBizVO.GroupId = item.GroupID
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
