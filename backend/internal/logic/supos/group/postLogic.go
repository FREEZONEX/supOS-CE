// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package group

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/utils/apiutil"
	"context"
	"time"
	"unicode/utf8"

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

func (l *PostLogic) Post(req *types.SaveGroupReq) (resp *types.JsonResult, err error) {

	name := *req.Name
	if req.Name == nil || name == "" {
		return nil, errors.Parameter.WithMsg("组名称不能为空")
	}
	db := stores.GetCommonConn(l.ctx)
	groupMapper := &relationDB.GroupMapper{}

	if utf8.RuneCountInString(name) > 255 {
		message := I18nUtils.GetMessage("group.name.maxLength")
		return nil, errors.Database.WithMsg(message).AddDetail(err)
	}

	var desc string
	if req.Description != nil {
		desc = *req.Description
	}

	if utf8.RuneCountInString(desc) > 255 {
		message := I18nUtils.GetMessage("group.name.maxLength")
		return nil, errors.Database.WithMsg(message).AddDetail(err)
	}

	// 检查组名是否已存在
	existingGroup, err := groupMapper.SelectByName(db, name)
	if err != nil {
		l.Errorf("查询组失败: %v", err)
		return nil, errors.Database.WithMsg("创建组失败").AddDetail(err)
	}
	if len(existingGroup) > 0 {
		return nil, errors.Parameter.WithMsg("组名称已存在")
	}

	// 构建组模型
	group := &relationDB.GroupModel{
		Type:        req.Type,
		Name:        name,
		Description: "",
		Sort:        0,
		UpdateAt:    time.Now(),
		CreateAt:    time.Now(),
	}

	// 设置创建者
	if userCtx := apiutil.GetUserFromContext(l.ctx); userCtx != nil {
		group.Creator = userCtx.PreferredUsername
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

	return &types.JsonResult{
		Code: 0,
		Msg:  "创建成功",
		Data: map[string]int64{"id": group.ID},
	}, nil
}
