// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package group

import (
	"context"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/stores"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据ID删除组
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.GroupIDReq) (*types.JsonResult, error) {
	db := stores.GetCommonConn(l.ctx)
	groupMapper := &relationDB.GroupMapper{}

	// 查询组是否存在
	group, err := groupMapper.SelectById(db, req.ID)
	if err != nil {
		l.Errorf("查询组失败: %v", err)
		return &types.JsonResult{
			Code: 500,
			Msg:  "查询组失败: " + err.Error(),
		}, nil
	}
	if group == nil {
		return &types.JsonResult{
			Code: 400,
			Msg:  "组不存在",
		}, nil
	}

	// 执行删除
	if err = groupMapper.DeleteById(db, req.ID); err != nil {
		l.Errorf("删除组失败: %v", err)
		return &types.JsonResult{
			Code: 500,
			Msg:  "删除组失败: " + err.Error(),
		}, nil
	}

	if err = groupMapper.DeleteById(db, req.ID); err != nil {
		l.Errorf("删除组关联的业务groupId失败: %v", err)
		return &types.JsonResult{
			Code: 500,
			Msg:  "删除组关联的业务groupId失败: " + err.Error(),
		}, nil
	}

	return &types.JsonResult{
		Code: 200,
		Msg:  "success",
		Data: nil,
	}, nil
}
