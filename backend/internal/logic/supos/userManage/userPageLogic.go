package userManage

import (
	"context"
	"strconv"
	"strings"

	"backend/internal/common/constants"
	"backend/internal/common/enums"
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

	db, err := l.keycloakDB()
	if err != nil {
		return nil, err
	}

	realm := strings.TrimSpace(l.realm())
	if realm == "" {
		return nil, errors.System.WithMsg("realm not configured")
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

	roleName := strings.TrimSpace(req.RoleName)
	if roleName != "" {
		if roleName == enums.RoleSuperAdmin.Comment || roleName == enums.RoleSuperAdmin.Name {
			roleName = enums.RoleSuperAdmin.Name
		}
	}

	base := db.Table("user_entity AS u").
		Joins("JOIN user_role_mapping AS urm ON urm.user_id = u.id").
		Joins("JOIN keycloak_role AS r ON r.id = urm.role_id").
		Where("u.realm_id = (SELECT id FROM realm WHERE name = ?)", realm).
		Where("u.service_account_client_link IS NULL")

	if v := strings.TrimSpace(req.PreferredUsername); v != "" {
		base = base.Where("u.username LIKE ?", "%"+v+"%")
	}
	if v := strings.TrimSpace(req.FirstName); v != "" {
		base = base.Where("u.first_name LIKE ?", "%"+v+"%")
	}
	if v := strings.TrimSpace(req.Email); v != "" {
		base = base.Where("u.email = ?", v)
	}
	if roleName != "" {
		base = base.Where("r.name = ?", roleName)
	}
	if req.Enabled != nil {
		base = base.Where("u.enabled = ?", *req.Enabled)
	}

	var total int64
	countQuery := base.Session(&gorm.Session{}).Distinct("u.id")
	if err := countQuery.Count(&total).Error; err != nil {
		l.Errorf("failed to count users: %v", err)
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

	var rows []userEntityRow
	dataQuery := base.
		Select(`u.id, u.username AS preferred_username, u.first_name, u.email, u.enabled, u.email_verified`).
		Group("u.id").
		Order("u.created_timestamp ASC").
		Limit(int(pageSize)).
		Offset(int(offset))

	if err := dataQuery.Scan(&rows).Error; err != nil {
		l.Errorf("failed to load user list: %v", err)
		return nil, errors.System.WithMsg("failed to query users")
	}

	userIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.ID)
	}

	attrMap, err := l.loadAttributesForUsers(db, userIDs)
	if err != nil {
		return nil, err
	}

	roleMap, err := l.loadRolesForUsers(db, userIDs)
	if err != nil {
		return nil, err
	}

	items := make([]types.UserManageItem, 0, len(rows))
	for _, row := range rows {
		item := types.UserManageItem{
			ID:                row.ID,
			PreferredUsername: row.PreferredUsername,
			FirstName:         row.FirstName,
			Email:             row.Email,
			Enabled:           row.Enabled,
			EmailVerified:     row.EmailVerified,
			Sub:               row.ID,
			FirstTimeLogin:    1,
			TipsEnable:        1,
			HomePage:          constants.DefaultHomepage,
		}

		if attrs := attrMap[row.ID]; len(attrs) > 0 {
			if v := strings.TrimSpace(attrs["phone"]); v != "" {
				item.Phone = v
			}
			if v := strings.TrimSpace(attrs["homePage"]); v != "" {
				item.HomePage = v
			}
			if v := strings.TrimSpace(attrs["source"]); v != "" {
				item.Source = v
			}
			if v := strings.TrimSpace(attrs["firstTimeLogin"]); v != "" {
				if iv, err := strconv.Atoi(v); err == nil {
					item.FirstTimeLogin = iv
				}
			}
			if v := strings.TrimSpace(attrs["tipsEnable"]); v != "" {
				if iv, err := strconv.Atoi(v); err == nil {
					item.TipsEnable = iv
				}
			}
		}

		if summaries := roleMap[row.ID]; len(summaries) > 0 {
			item.RoleList = summaries
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

type userEntityRow struct {
	ID                string `gorm:"column:id"`
	PreferredUsername string `gorm:"column:preferred_username"`
	FirstName         string `gorm:"column:first_name"`
	Email             string `gorm:"column:email"`
	Enabled           bool   `gorm:"column:enabled"`
	EmailVerified     bool   `gorm:"column:email_verified"`
}
