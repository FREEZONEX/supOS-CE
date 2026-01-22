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

type PostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建组
func NewPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostLogic {
	return &PostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostLogic) Post(req *types.SaveGroupReq) (resp *types.CommonResp, err error) {
	if req.Name == nil || *req.Name == "" {
		return nil, errors.Parameter.WithMsg("组名称不能为空")
	}

	db := stores.GetCommonConn(l.ctx)
	groupMapper := &relationDB.GroupMapper{}

	// 构建组模型
	group := &relationDB.GroupModel{
		Type:        req.Type,
		Name:        *req.Name,
		Description: "",
		Sort:        0,
		UpdateAt:    time.Now(),
		CreateAt:    time.Now(),
	}
	if req.Description != nil {
		group.Description = *req.Description
	}
	if req.Sort != nil {
		group.Sort = *req.Sort
	}

	// 插入数据
	if err = groupMapper.Insert(db, group); err != nil {
		l.Errorf("创建组失败: %v", err)
		return nil, errors.Database.WithMsg("创建组失败").AddDetail(err)
	}

	return &types.CommonResp{
		ID: group.ID,
	}, nil
}
