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

type BatchDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除组
func NewBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteLogic {
	return &BatchDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteLogic) BatchDelete(req *types.GroupBatchDeleteReq) error {
	if len(req.IDs) == 0 {
		return errors.Parameter.WithMsg("请选择要删除的组")
	}

	db := stores.GetCommonConn(l.ctx)
	groupMapper := &relationDB.GroupMapper{}

	// 查询组是否存在
	groups, err := groupMapper.SelectByIds(db, req.IDs)
	if err != nil {
		l.Errorf("查询组列表失败: %v", err)
		return errors.Database.WithMsg("查询组列表失败").AddDetail(err)
	}

	// 检查是否所有ID都存在
	if len(groups) != len(req.IDs) {
		return errors.Parameter.WithMsg("部分组不存在")
	}

	// 执行批量删除
	if err = groupMapper.DeleteBatchIds(db, req.IDs); err != nil {
		l.Errorf("批量删除组失败: %v", err)
		return errors.Database.WithMsg("批量删除组失败").AddDetail(err)
	}

	return nil
}
