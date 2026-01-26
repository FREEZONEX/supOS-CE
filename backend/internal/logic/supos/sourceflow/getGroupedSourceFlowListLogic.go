// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"backend/internal/repo/relationDB"
	"context"
	"fmt"
	"net/http"
	"time"

	"backend/internal/svc"
	"backend/internal/types"

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
	req.GroupType = 1
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
			groupBizVO.CreateAt = item.CreateAt.Format(time.RFC3339)
			groupBizVO.GroupId = item.GroupID
		} else if item.GroupType == 0 {
			// 处理未分组的source flow数据
			bizID, _ := parseStringToInt64(item.ID)
			groupBizVO.ID = bizID
			groupBizVO.Name = item.Name
			groupBizVO.Description = item.Description
			groupBizVO.CreateAt = item.CreateAt.Format(time.RFC3339)
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
