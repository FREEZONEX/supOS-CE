// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package group

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"context"
	"time"
	"unicode/utf8"

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
	ctx := l.ctx
	if req.ID == nil {
		return nil, errors.Parameter.WithMsg(I18nUtils.GetMessageWithCtx(ctx, "group.id.isnotblank"))
	}

	db := stores.GetCommonConn(l.ctx)
	groupMapper := &relationDB.GroupMapper{}

	// 查询组是否存在
	group, err := groupMapper.SelectById(db, *req.ID)
	if err != nil {
		l.Errorf("查询组失败: %v", err)
		return nil, errors.Database.WithMsg(I18nUtils.GetMessageWithCtx(ctx, "group.query.failed")).AddDetail(err)
	}
	if group == nil {
		return nil, errors.Parameter.WithMsg(I18nUtils.GetMessageWithCtx(ctx, "group.notfound"))
	}

	if utf8.RuneCountInString(*req.Name) > 255 {
		message := I18nUtils.GetMessageWithCtx(ctx, "group.name.maxLength")
		return nil, errors.Database.WithMsg(message).AddDetail(err)
	}

	var desc string
	if req.Description != nil {
		desc = *req.Description
	}

	if utf8.RuneCountInString(desc) > 512 {
		message := I18nUtils.GetMessageWithCtx(ctx, "group.description.maxLength")
		return nil, errors.Database.WithMsg(message).AddDetail(err)
	}

	// 检查组名是否已存在
	existingGroup, err := groupMapper.SelectByNameNotId(db, *req.ID, *req.Name, *req.Type)
	if err != nil {
		l.Errorf("查询组失败: %v", err)
		return nil, errors.Database.WithMsg(I18nUtils.GetMessageWithCtx(ctx, "group.update.failed")).AddDetail(err)
	}
	if len(existingGroup) > 0 {
		return nil, errors.Parameter.WithMsg(I18nUtils.GetMessageWithCtx(ctx, "group.name.duplication"))
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
		return nil, errors.Database.WithMsg(I18nUtils.GetMessageWithCtx(ctx, "group.update.failed")).AddDetail(err)
	}

	return &types.JsonResult{
		Code: 0,
		Msg:  I18nUtils.GetMessageWithCtx(ctx, "group.update.success"),
		Data: nil,
	}, nil
}
