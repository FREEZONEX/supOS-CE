// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"
	"strings"
	"time"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMembersLogic {
	return &GetMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMembersLogic) GetMembers(req *types.GetMembersReq) (*types.GetMembersResp, error) {
	if req == nil {
		req = &types.GetMembersReq{}
	}
	pageNo, pageSize := normalizeGetMembersPage(req.PageNo, req.PageSize)

	db := stores.GetCommonConn(l.ctx)
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}

	query := db.WithContext(l.ctx).Table("supos_user AS u")
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, err
	}

	type userRow struct {
		ID        string    `gorm:"column:id"`
		Username  string    `gorm:"column:username"`
		Email     string    `gorm:"column:email"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}
	var rows []userRow
	if err := query.Session(&gorm.Session{}).
		Select("u.id, u.username, u.email, u.updated_at").
		Order("u.updated_at DESC, u.id ASC").
		Limit(pageSize).
		Offset((pageNo - 1) * pageSize).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	list := make([]types.MemberBrief, 0, len(rows))
	for _, row := range rows {
		list = append(list, types.MemberBrief{
			UserID:    row.ID,
			UserName:  row.Username,
			Email:     strings.TrimSpace(row.Email),
			UpdatedAt: row.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &types.GetMembersResp{
		PageNo:   pageNo,
		PageSize: pageSize,
		Total:    int(total),
		List:     list,
	}, nil
}

func normalizeGetMembersPage(pageNo, pageSize int) (int, int) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageNo, pageSize
}
