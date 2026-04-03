package userManage

import (
	"context"
	"strings"

	"backend/internal/common/constants"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gorm.io/gorm"
)

type UserPageLogic struct {
	baseUserManageLogic
}

// Query paginated user list
func NewUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPageLogic {
	return &UserPageLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *UserPageLogic) UserPage(req *types.UserManagePageReq) (*types.UserManagePageResp, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("request body is empty")
	}

	db, err := l.db()
	if err != nil {
		return nil, err
	}

	pageNo := req.PageNo
	if pageNo <= 0 {
		pageNo = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	query := db.Table("supos_user AS u")
	roleName := strings.TrimSpace(req.RoleName)
	if roleName != "" {
		query = query.Joins("JOIN supos_user_role AS ur ON ur.user_id = u.id").
			Joins("JOIN supos_role AS r ON r.id = ur.role_id AND r.status = 1")
		if builtinKey := normalizeBuiltinRoleKey(roleName); builtinKey != "" {
			query = query.Where("LOWER(r.role_key) = LOWER(?)", builtinKey)
		} else {
			query = query.Where("LOWER(r.role_name) = LOWER(?) OR LOWER(r.role_key) = LOWER(?)", roleName, roleName)
		}
	}

	if v := strings.TrimSpace(req.PreferredUsername); v != "" {
		query = query.Where("u.username ILIKE ?", "%"+v+"%")
	}
	if v := strings.TrimSpace(req.FirstName); v != "" {
		query = query.Where("u.display_name ILIKE ?", "%"+v+"%")
	}
	if v := strings.TrimSpace(req.Email); v != "" {
		query = query.Where("LOWER(u.email) = LOWER(?)", v)
	}
	if req.Enabled != nil {
		query = query.Where("u.enabled = ?", *req.Enabled)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Distinct("u.id").Count(&total).Error; err != nil {
		l.Errorf("count users failed: %v", err)
		return nil, errors.System.WithMsg("failed to query users")
	}
	if total == 0 {
		return &types.UserManagePageResp{
			PageNo:   pageNo,
			PageSize: pageSize,
			Total:    0,
			Data:     nil,
		}, nil
	}

	var rows []relationDB.IamUser
	if err := query.Session(&gorm.Session{}).
		Select("u.*").
		Distinct("u.id, u.username, u.display_name, u.email, u.enabled, u.source, u.phone, u.home_page, u.first_time_login, u.tips_enable, u.main_language, u.created_at, u.updated_at").
		Order("u.created_at ASC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(&rows).Error; err != nil {
		l.Errorf("load users failed: %v", err)
		return nil, errors.System.WithMsg("failed to query users")
	}

	userIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.ID)
	}

	roleMap, err := l.loadRolesForUsers(db, userIDs)
	if err != nil {
		return nil, err
	}

	items := make([]types.UserManageItem, 0, len(rows))
	for _, row := range rows {
		item := types.UserManageItem{
			ID:                row.ID,
			Email:             strings.TrimSpace(row.Email),
			EmailVerified:     false,
			FirstName:         firstNonEmpty(row.DisplayName, row.Username),
			PreferredUsername: row.Username,
			Sub:               row.ID,
			Enabled:           row.Enabled,
			FirstTimeLogin:    row.FirstTimeLogin,
			TipsEnable:        row.TipsEnable,
			HomePage:          firstNonEmpty(row.HomePage, constants.DefaultHomepage),
			Phone:             strings.TrimSpace(row.Phone),
			Source:            strings.TrimSpace(row.Source),
			RoleList:          roleMap[row.ID],
		}
		if item.HomePage == "" {
			item.HomePage = constants.DefaultHomepage
		}
		items = append(items, item)
	}

	return &types.UserManagePageResp{
		Code:     0,
		PageNo:   pageNo,
		PageSize: pageSize,
		Total:    total,
		Data:     items,
	}, nil
}
