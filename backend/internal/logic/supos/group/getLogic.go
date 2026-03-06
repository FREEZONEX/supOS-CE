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

type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询组列表
func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogic) Get(req *types.GroupQuery) (resp []types.GroupVO, err error) {
	db := stores.GetCommonConn(l.ctx)
	query := db.WithContext(l.ctx).Model(&relationDB.GroupModel{})

	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.Name != nil && *req.Name != "" {
		query = query.Where("name LIKE ?", "%"+*req.Name+"%")
	}

	var total int64
	if err = query.Count(&total).Error; err != nil {
		l.Errorf("查询组总数失败: %v", err)
		return nil, errors.Database.WithMsg("查询组列表失败").AddDetail(err)
	}

	var groups []*relationDB.GroupModel
	if err = query.Order("sort ASC, id ASC").
		Limit(int(req.Size)).Offset(int((req.Page - 1) * req.Size)).
		Find(&groups).Error; err != nil {
		l.Errorf("查询组列表失败: %v", err)
		return nil, errors.Database.WithMsg("查询组列表失败").AddDetail(err)
	}

	resp = make([]types.GroupVO, 0, len(groups))
	for _, group := range groups {
		resp = append(resp, types.GroupVO{
			ID:          group.ID,
			Type:        group.Type,
			Name:        group.Name,
			Description: group.Description,
			Sort:        group.Sort,
			UpdateAt:    group.UpdateAt.Format(time.RFC3339),
			CreateAt:    group.CreateAt.Format(time.RFC3339),
		})
	}

	return resp, nil
}
