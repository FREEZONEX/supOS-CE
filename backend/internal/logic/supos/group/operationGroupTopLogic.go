// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package group

import (
	"backend/internal/repo/relationDB"
	"context"
	"net/http"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationGroupTopLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	groupMapper relationDB.GroupMapper
}

// 置顶分组
func NewOperationGroupTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationGroupTopLogic {
	return &OperationGroupTopLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		groupMapper: relationDB.GroupMapper{},
	}
}

func (l *OperationGroupTopLogic) OperationGroupTop(req *types.OperationGroupTopReq) (resp *types.JsonResult, err error) {
	db := relationDB.GetDb(l.ctx)

	// 根据 status 设置 sort 值
	var sort int32
	if *req.Status {
		// 置顶：设置 sort 为 1
		sort = 1
	} else {
		// 取消置顶：设置 sort 为 0
		sort = 0
	}

	// 更新 resource_group 表的 sort 字段
	err = l.groupMapper.UpdateSortById(db, *req.ID, sort)
	if err != nil {
		l.Logger.Error("更新分组排序失败:", err)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  "更新分组排序失败",
			Data: nil,
		}, nil
	}

	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "操作成功",
		Data: nil,
	}, nil
}
