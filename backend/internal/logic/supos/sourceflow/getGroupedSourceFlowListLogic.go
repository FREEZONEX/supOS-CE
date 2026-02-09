// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupedSourceFlowListLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	sourceFlowRepo *relationDB.NoderedSourceFlowRepo
}

// 分页按分组获取source flow列表
func NewGetGroupedSourceFlowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupedSourceFlowListLogic {
	return &GetGroupedSourceFlowListLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		sourceFlowRepo: relationDB.NewNoderedSourceFlowRepo(ctx),
	}
}

func (l *GetGroupedSourceFlowListLogic) GetGroupedSourceFlowList(req *types.GroupPageRequest) (resp *types.SourceFlowGroupPageResp, err error) {
	if req.PageNo < 1 {
		req.PageNo = 1
	}
	pageResult := &types.SourceFlowGroupPageResp{
		PageResultDTO: types.PageResultDTO{
			Code:     http.StatusOK,
			PageNo:   req.PageNo,
			PageSize: req.PageSize,
		},
	}

	l.Logger.Debugf("GetGroupedSourceFlowList: PageList request: %+v", req)
	db := relationDB.GetDb(l.ctx)

	// 调用DAO层方法获取分组后的source flow列表
	items, total, err := l.sourceFlowRepo.GetGroupedSourceFlowList(db, req)
	if err != nil {
		l.Logger.Error("查询分组source flow列表失败:", err)
		return nil, err
	}

	// 转换为前端需要的格式
	var groupVOList []types.GroupFlowVO
	for _, item := range items {
		groupFlowVO := types.GroupFlowVO{}
		//存在分组
		groupFlowVO.ID = item.ID
		groupFlowVO.Category = item.Category
		groupFlowVO.Name = item.Name
		groupFlowVO.Description = item.Description
		groupFlowVO.Sort = item.Sort
		groupFlowVO.CreateAt = item.CreateAt.UnixMilli()
		groupFlowVO.Creator = item.Creator
		groupFlowVO.FlowName = item.FlowName
		groupFlowVO.FlowID = item.FlowID
		groupFlowVO.FlowStatus = item.FlowStatus
		groupFlowVO.Template = item.Template
		groupFlowVO.HasChildren = item.HasChildren
		if item.GroupType != 0 {
			groupFlowVO.GroupType = &item.GroupType
		}
		groupVOList = append(groupVOList, groupFlowVO)
	}

	pageResult.Total = total
	pageResult.Data = groupVOList

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
