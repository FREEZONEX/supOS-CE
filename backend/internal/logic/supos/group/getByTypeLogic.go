// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package group

import (
	"context"
	"time"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据类型查询组列表
func NewGetByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByTypeLogic {
	return &GetByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetByTypeLogic) GetByType(req *types.GroupByTypeQuery) (resp *types.GroupPageResultVO, err error) {
	db := stores.GetCommonConn(l.ctx)

	// 使用 groupDao 查询
	groupMapper := &relationDB.GroupMapper{}
	groups, err := groupMapper.SelectByType(db, req.Type)
	if err != nil {
		l.Errorf("根据类型查询组列表失败: %v", err)
		return nil, errors.Database.WithMsg("查询组列表失败").AddDetail(err)
	}

	data := make([]types.GroupVO, 0, len(groups))
	for _, group := range groups {
		data = append(data, types.GroupVO{
			ID:          group.ID,
			Type:        group.Type,
			Name:        group.Name,
			Description: group.Description,
			Sort:        group.Sort,
			UpdateAt:    group.UpdateAt.Format(time.RFC3339),
			CreateAt:    group.CreateAt.Format(time.RFC3339),
		})
	}

	// 使用分页对象包装结果
	// 注意：这里 SelectByType 查询的是所有数据，我们需要手动分页
	pageNo := req.Page
	if pageNo <= 0 {
		pageNo = 1
	}
	pageSize := req.Size
	if pageSize <= 0 {
		pageSize = 20
	}

	total := int64(len(data))
	start := (pageNo - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedData := data[start:end]

	return &types.GroupPageResultVO{
		PageNo:   pageNo,
		PageSize: pageSize,
		Total:    total,
		Code:     200,
		Data:     pagedData,
	}, nil
}
