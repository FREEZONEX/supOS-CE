// Code scaffolded by goctl. Safe to edit.

package openservice

import (
	"context"
	"strings"

	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserOpenapiService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserOpenapiService creates a new UserOpenapiService instance
func NewUserOpenapiService(ctx context.Context, svcCtx *svc.ServiceContext) *UserOpenapiService {
	return &UserOpenapiService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserPageQueryDto 用户分页查询参数
type UserPageQueryDto struct {
	PageNo        int
	PageSize      int
	Username      string
	DisplayName   string
	Email         string
	Phone         string
	Enabled       *bool
	ExactUsername string
}

// UserManageResult 用户管理查询结果
type UserManageResult struct {
	Users []types.OpenUserInfo
	Total int
}

type openUserRow struct {
	ID             string `gorm:"column:id"`
	Username       string `gorm:"column:username"`
	DisplayName    string `gorm:"column:display_name"`
	Email          string `gorm:"column:email"`
	Enabled        bool   `gorm:"column:enabled"`
	Source         string `gorm:"column:source"`
	Phone          string `gorm:"column:phone"`
	HomePage       string `gorm:"column:home_page"`
	FirstTimeLogin int    `gorm:"column:first_time_login"`
	TipsEnable     int    `gorm:"column:tips_enable"`
}

// UserManageList 获取用户列表（供列表和详情接口共用）
func (s *UserOpenapiService) UserManageList(params UserPageQueryDto) (*UserManageResult, error) {
	db := stores.GetCommonConn(s.ctx)
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	query := db.WithContext(s.ctx).
		Table("supos_user AS u")

	if exactUsername := strings.TrimSpace(params.ExactUsername); exactUsername != "" {
		query = query.Where("LOWER(u.username) = LOWER(?)", exactUsername)
	} else {
		if username := strings.TrimSpace(params.Username); username != "" {
			query = query.Where("LOWER(u.username) = LOWER(?)", username)
		}
		if displayName := strings.TrimSpace(params.DisplayName); displayName != "" {
			query = query.Where("u.display_name ILIKE ?", "%"+displayName+"%")
		}
		if email := strings.TrimSpace(params.Email); email != "" {
			query = query.Where("u.email ILIKE ?", "%"+email+"%")
		}
		if phone := strings.TrimSpace(params.Phone); phone != "" {
			query = query.Where("u.phone ILIKE ?", "%"+phone+"%")
		}
	}
	if params.Enabled != nil {
		query = query.Where("u.enabled = ?", *params.Enabled)
	}

	pageNo := params.PageNo
	if pageNo <= 0 {
		pageNo = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (pageNo - 1) * pageSize

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		s.Errorf("Failed to count users: %v", err)
		return nil, err
	}

	var rows []openUserRow
	if err := query.Session(&gorm.Session{}).
		Select("u.id, u.username, u.display_name, u.email, u.enabled, u.source, u.phone, u.home_page, u.first_time_login, u.tips_enable").
		Order("u.created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Scan(&rows).Error; err != nil {
		s.Errorf("Failed to query users: %v", err)
		return nil, err
	}

	roleMap, err := loadOpenUserRoles(s.ctx, rows)
	if err != nil {
		s.Errorf("Failed to load user roles: %v", err)
		return nil, err
	}

	resultList := make([]types.OpenUserInfo, 0, len(rows))
	for _, row := range rows {
		homePage := strings.TrimSpace(row.HomePage)
		if homePage == "" {
			homePage = constants.DefaultHomepage
		}
		firstName := strings.TrimSpace(row.DisplayName)
		if firstName == "" {
			firstName = row.Username
		}
		firstTimeLogin := row.FirstTimeLogin
		if firstTimeLogin == 0 {
			firstTimeLogin = 1
		}
		tipsEnable := row.TipsEnable
		if tipsEnable == 0 {
			tipsEnable = 1
		}
		resultList = append(resultList, types.OpenUserInfo{
			ID:                row.ID,
			Email:             strings.TrimSpace(row.Email),
			EmailVerified:     false,
			FirstName:         firstName,
			PreferredUsername: row.Username,
			Sub:               row.ID,
			Enabled:           row.Enabled,
			RoleList:          roleMap[row.ID],
			FirstTimeLogin:    firstTimeLogin,
			TipsEnable:        tipsEnable,
			HomePage:          homePage,
			Phone:             strings.TrimSpace(row.Phone),
			Source:            strings.TrimSpace(row.Source),
		})
	}

	return &UserManageResult{
		Users: resultList,
		Total: int(total),
	}, nil
}

func loadOpenUserRoles(ctx context.Context, users []openUserRow) (map[string][]types.RoleSummary, error) {
	if len(users) == 0 {
		return map[string][]types.RoleSummary{}, nil
	}
	db := stores.GetCommonConn(ctx)
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}

	userIDs := make([]string, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	type roleJoinRow struct {
		UserID      string `gorm:"column:user_id"`
		RoleID      string `gorm:"column:role_id"`
		RoleKey     string `gorm:"column:role_key"`
		RoleName    string `gorm:"column:role_name"`
		Description string `gorm:"column:description"`
	}
	var rows []roleJoinRow
	if err := db.WithContext(ctx).
		Table("supos_user_role AS ur").
		Select("ur.user_id, r.id AS role_id, r.role_key, r.role_name, r.description").
		Joins("JOIN supos_role AS r ON r.id = ur.role_id").
		Where("ur.user_id IN ?", userIDs).
		Where("r.status = 1").
		Order("r.builtin DESC, r.role_name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string][]types.RoleSummary)
	for _, row := range rows {
		display := strings.TrimSpace(row.RoleName)
		switch normalizeOpenRoleKey(firstNonEmptyOpen(row.RoleKey, row.RoleName)) {
		case "admin":
			display = I18nUtils.GetMessageWithCtx(ctx, enums.RoleAdmin.I18nCode)
		case "user":
			display = I18nUtils.GetMessageWithCtx(ctx, enums.RoleNormalUser.I18nCode)
		}
		result[row.UserID] = append(result[row.UserID], types.RoleSummary{
			RoleID:          row.RoleID,
			RoleName:        firstNonEmptyOpen(display, row.RoleName, row.RoleKey),
			RoleDescription: strings.TrimSpace(row.Description),
			ClientRole:      true,
		})
	}
	return result, nil
}

func normalizeOpenRoleKey(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "admin", "super-admin", strings.ToLower(enums.RoleAdmin.Comment), strings.ToLower(enums.RoleSuperAdmin.Comment):
		return "admin"
	case "user", "normal-user", strings.ToLower(enums.RoleNormalUser.Comment):
		return "user"
	default:
		return ""
	}
}

func firstNonEmptyOpen(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
