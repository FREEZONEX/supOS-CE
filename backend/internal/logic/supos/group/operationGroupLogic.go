// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package group

import (
	"backend/internal/repo/relationDB"
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type OperationGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 操作分组数据（移入移出）
func NewOperationGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationGroupLogic {
	return &OperationGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 操作分组数据（移入移出） 根据 req 中的参数执行相应的操作
func (l *OperationGroupLogic) OperationGroup(req *types.OperationGroupReq) (resp *types.JsonResult, err error) {
	groupMapper := &relationDB.GroupMapper{}
	err = groupMapper.OperationGroup(relationDB.GetDb(l.ctx), req)
	if err != nil {
		logx.Errorf("操作分组失败: %v", err)
		switch err {
		case errors.NotFind:
			return &types.JsonResult{
				Code: 404,
				Msg:  "分组不存在",
			}, nil
		case errors.Parameter:
			return &types.JsonResult{
				Code: 400,
				Msg:  "不支持的分组类型",
			}, nil
		default:
			return &types.JsonResult{
				Code: 500,
				Msg:  "操作分组失败",
			}, nil
		}
	}

	return &types.JsonResult{
		Code: 200,
		Msg:  "操作成功",
		Data: &types.OperationResult{Success: true},
	}, nil
}
