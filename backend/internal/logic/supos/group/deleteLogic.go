// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package group

import (
	"context"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/stores"

	"gitee.com/unitedrhino/share/errors"
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

func (l *DeleteLogic) Delete(req *types.GroupIDReq) error {
	db := stores.GetCommonConn(l.ctx)
	groupMapper := &relationDB.GroupMapper{}

	// 查询组是否存在
	group, err := groupMapper.SelectById(db, req.ID)
	if err != nil {
		l.Errorf("查询组失败: %v", err)
		return errors.Database.WithMsg("查询组失败").AddDetail(err)
	}
	if group == nil {
		return errors.Parameter.WithMsg("组不存在")
	}

	// 执行删除
	if err = groupMapper.DeleteById(db, req.ID); err != nil {
		l.Errorf("删除组失败: %v", err)
		return errors.Database.WithMsg("删除组失败").AddDetail(err)
	}

	return nil
}
