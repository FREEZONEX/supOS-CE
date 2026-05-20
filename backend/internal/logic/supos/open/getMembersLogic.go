// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"
	"strings"
	"time"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
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
	updatedAtRange, err := parseGetMembersUpdatedAtRange(req)
	if err != nil {
		return nil, err
	}

	db := stores.GetCommonConn(l.ctx)
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}

	query := db.WithContext(l.ctx).Table("supos_user AS u")
	if updatedAtRange.Start != nil {
		query = query.Where("u.updated_at >= ?", *updatedAtRange.Start)
	}
	if updatedAtRange.End != nil {
		if updatedAtRange.EndExclusive {
			query = query.Where("u.updated_at < ?", *updatedAtRange.End)
		} else {
			query = query.Where("u.updated_at <= ?", *updatedAtRange.End)
		}
	}

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

type getMembersUpdatedAtRange struct {
	Start        *time.Time
	End          *time.Time
	EndExclusive bool
}

func parseGetMembersUpdatedAtRange(req *types.GetMembersReq) (getMembersUpdatedAtRange, error) {
	var result getMembersUpdatedAtRange
	if req == nil {
		return result, nil
	}
	if value := strings.TrimSpace(req.UpdatedAtStart); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return result, errors.Parameter.WithMsg("updatedAtStart.invalid")
		}
		result.Start = &parsed
	}

	if value := strings.TrimSpace(req.UpdatedAtEnd); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return result, errors.Parameter.WithMsg("updatedAtEnd.invalid")
		}
		if !rfc3339HasFractionalSecond(value) {
			parsed = parsed.Add(time.Second)
			result.EndExclusive = true
		}
		result.End = &parsed
	}
	return result, nil
}

func rfc3339HasFractionalSecond(value string) bool {
	tIndex := strings.IndexByte(value, 'T')
	if tIndex < 0 {
		return false
	}
	timePart := value[tIndex+1:]
	if zIndex := strings.IndexByte(timePart, 'Z'); zIndex >= 0 {
		timePart = timePart[:zIndex]
	} else {
		for i := len(timePart) - 1; i >= 0; i-- {
			if timePart[i] == '+' || timePart[i] == '-' {
				timePart = timePart[:i]
				break
			}
		}
	}
	return strings.Contains(timePart, ".")
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
