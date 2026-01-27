// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package group

import (
	"context"
	"time"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/stores"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type PutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新组
func NewPutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutLogic {
	return &PutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutLogic) Put(req *types.SaveGroupReq) (resp *types.JsonResult, err error) {
	if req.ID == nil {
		return nil, errors.Parameter.WithMsg("组ID不能为空")
	}

	db := stores.GetCommonConn(l.ctx)
	groupMapper := &relationDB.GroupMapper{}

	// 查询组是否存在
	group, err := groupMapper.SelectById(db, *req.ID)
	if err != nil {
		l.Errorf("查询组失败: %v", err)
		return nil, errors.Database.WithMsg("查询组失败").AddDetail(err)
	}
	if group == nil {
		return nil, errors.Parameter.WithMsg("组不存在")
	}

	// 更新字段
	if req.Type != nil {
		group.Type = req.Type
	}
	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.Description != nil {
		group.Description = *req.Description
	}
	if req.Sort != nil {
		group.Sort = *req.Sort
	}
	group.UpdateAt = time.Now()

	// 执行更新
	if err = groupMapper.UpdateById(db, group); err != nil {
		l.Errorf("更新组失败: %v", err)
		return nil, errors.Database.WithMsg("更新组失败").AddDetail(err)
	}

	return &types.JsonResult{
		Code: 0,
		Msg:  "更新成功",
		Data: nil,
	}, nil
}
