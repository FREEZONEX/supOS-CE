package iam

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/enums"
	"backend/internal/common/utils/langutil"
	"backend/internal/common/vo"
	"backend/internal/repo/relationDB"

	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepo struct {
	db *gorm.DB
}

type RoleSyncInput struct {
	Role      *authdto.RoleDto
	Resources []*authdto.ResourceDto
}

func NewAuthRepo(ctx context.Context) (*AuthRepo, error) {
	conn := stores.GetCommonConn(ctx)
	if conn == nil {
		return nil, stores.ErrFmt(fmt.Errorf("common database connection not initialized"))
	}
	return &AuthRepo{db: conn.WithContext(ctx)}, nil
}

func (r *AuthRepo) GetUserByID(userID string) (*relationDB.IamUser, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, nil
	}
	var user relationDB.IamUser
	err := r.db.Where("id = ?", strings.TrimSpace(userID)).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &user, nil
}

func (r *AuthRepo) GetUserByUsername(username string) (*relationDB.IamUser, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, nil
	}
	var user relationDB.IamUser
	err := r.db.Where("LOWER(username) = LOWER(?)", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &user, nil
}

func (r *AuthRepo) GetSession(sessionID string) (*relationDB.IamSession, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, nil
	}
	var session relationDB.IamSession
	err := r.db.Where("id = ?", sessionID).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &session, nil
}

func (r *AuthRepo) CreateSession(session *relationDB.IamSession) error {
	if session == nil {
		return nil
	}
	return stores.ErrFmt(r.db.Create(session).Error)
}

func (r *AuthRepo) TouchSession(sessionID string, now, expiredAt time.Time) error {
	return stores.ErrFmt(r.db.Model(&relationDB.IamSession{}).
		Where("id = ? AND revoked_at IS NULL", strings.TrimSpace(sessionID)).
		Updates(map[string]any{
			"last_access_at": now,
			"expired_at":     expiredAt,
		}).Error)
}

func (r *AuthRepo) RevokeSession(sessionID string, when time.Time) error {
	return stores.ErrFmt(r.db.Model(&relationDB.IamSession{}).
		Where("id = ?", strings.TrimSpace(sessionID)).
		Update("revoked_at", when).Error)
}

func (r *AuthRepo) BuildUserInfo(ctx context.Context, userID string, defaultHome string) (*vo.UserInfoVo, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, nil
	}

	userEntity, err := r.GetUserByID(userID)
	if err != nil || userEntity == nil {
		return nil, err
	}

	user := vo.NewUserInfoVo(userEntity.ID, userEntity.Username)
	user.Email = strings.TrimSpace(userEntity.Email)
	user.Enabled = userEntity.Enabled
	user.FirstName = strings.TrimSpace(userEntity.DisplayName)
	user.Source = strings.TrimSpace(userEntity.Source)
	user.Phone = strings.TrimSpace(userEntity.Phone)
	user.FirstTimeLogin = userEntity.FirstTimeLogin
	user.TipsEnable = userEntity.TipsEnable
	user.MainLanguage = strings.TrimSpace(userEntity.MainLanguage)
	user.HomePage = defaultHome

	if homePage := strings.TrimSpace(userEntity.HomePage); homePage != "" {
		user.HomePage = homePage
	}
	if user.MainLanguage == "" {
		user.MainLanguage = langutil.SystemLocale()
	}

	roles, err := r.getRolesByUserID(userID)
	if err != nil {
		return nil, err
	}
	user.RoleList = roles

	pages, buttons, err := r.getPermissionsByRoleIDs(ctx, collectRoleIDs(roles))
	if err != nil {
		return nil, err
	}
	user.ResourceList = appendDefaultResources(pages)
	user.ButtonList = buttons
	user.SuperAdmin = user.IsSuperAdmin()
	return user, nil
}

func (r *AuthRepo) SyncUser(ctx context.Context, user *vo.UserInfoVo, passwordHash string, roleInputs []RoleSyncInput) error {
	if user == nil || strings.TrimSpace(user.Sub) == "" || strings.TrimSpace(user.PreferredUsername) == "" {
		return nil
	}
	return stores.ErrFmt(r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		userModel := relationDB.IamUser{
			ID:             strings.TrimSpace(user.Sub),
			Username:       strings.TrimSpace(user.PreferredUsername),
			DisplayName:    strings.TrimSpace(firstNonEmpty(user.FirstName, user.PreferredUsername)),
			Email:          strings.TrimSpace(user.Email),
			Enabled:        user.Enabled,
			Source:         strings.TrimSpace(user.Source),
			Phone:          strings.TrimSpace(user.Phone),
			HomePage:       strings.TrimSpace(user.HomePage),
			FirstTimeLogin: user.FirstTimeLogin,
			TipsEnable:     user.TipsEnable,
			MainLanguage:   strings.TrimSpace(user.MainLanguage),
			UpdatedAt:      now,
		}
		updateColumns := []string{
			"username", "display_name", "email", "enabled", "source",
			"phone", "home_page", "first_time_login", "tips_enable", "main_language", "updated_at",
		}
		if passwordHash != "" {
			userModel.Password = passwordHash
			updateColumns = append(updateColumns, "password")
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns(updateColumns),
		}).Create(&userModel).Error; err != nil {
			return err
		}

		normalizedRoles := normalizeRoleInputs(roleInputs)
		roleIDs := make([]string, 0, len(normalizedRoles))
		for _, input := range normalizedRoles {
			if input.Role == nil {
				continue
			}
			role := relationDB.IamRole{
				ID:          strings.TrimSpace(input.Role.RoleID),
				RoleKey:     normalizeRoleKey(input.Role),
				RoleName:    strings.TrimSpace(input.Role.RoleName),
				Description: strings.TrimSpace(input.Role.RoleDescription),
				Builtin:     isBuiltinRole(input.Role),
				Status:      1,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"role_key", "role_name", "description", "builtin", "status", "updated_at"}),
			}).Create(&role).Error; err != nil {
				return err
			}
			roleIDs = append(roleIDs, role.ID)

			resourceIDs, err := mapResourcesToIDs(tx, input.Resources)
			if err != nil {
				return err
			}
			if err := tx.Where("role_id = ?", role.ID).Delete(&relationDB.IamRoleResource{}).Error; err != nil {
				return err
			}
			if len(resourceIDs) == 0 {
				continue
			}
			roleResources := make([]relationDB.IamRoleResource, 0, len(resourceIDs))
			for _, resourceID := range resourceIDs {
				roleResources = append(roleResources, relationDB.IamRoleResource{
					RoleID:     role.ID,
					ResourceID: resourceID,
					CreatedAt:  now,
				})
			}
			if err := tx.Create(&roleResources).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("user_id = ?", userModel.ID).Delete(&relationDB.IamUserRole{}).Error; err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			roleIDs = []string{enums.RoleNormalUser.ID}
		}

		assignments := make([]relationDB.IamUserRole, 0, len(roleIDs))
		for _, roleID := range uniqueStrings(roleIDs) {
			assignments = append(assignments, relationDB.IamUserRole{
				UserID:    userModel.ID,
				RoleID:    roleID,
				CreatedAt: now,
			})
		}
		return tx.Create(&assignments).Error
	}))
}

func (r *AuthRepo) getRolesByUserID(userID string) ([]*authdto.RoleDto, error) {
	var rows []relationDB.IamRole
	err := r.db.Table("supos_role AS r").
		Select("r.*").
		Joins("JOIN supos_user_role AS ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = 1", userID).
		Order("r.builtin DESC, r.role_name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	result := make([]*authdto.RoleDto, 0, len(rows))
	for _, row := range rows {
		result = append(result, &authdto.RoleDto{
			RoleID:          row.ID,
			RoleName:        row.RoleName,
			RoleDescription: row.Description,
			ClientRole:      true,
		})
	}
	return result, nil
}

func (r *AuthRepo) getPermissionsByRoleIDs(ctx context.Context, roleIDs []string) ([]*authdto.ResourceDto, []string, error) {
	if len(roleIDs) == 0 {
		return nil, nil, nil
	}
	var rows []relationDB.SuposResource
	err := r.db.WithContext(ctx).
		Table("supos_resource AS sr").
		Select("sr.*").
		Joins("JOIN supos_role_resource AS rr ON rr.resource_id = sr.id").
		Where("rr.role_id IN ?", uniqueStrings(roleIDs)).
		Where("COALESCE(sr.enable, true) = true").
		Order("sr.sort ASC, sr.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, nil, stores.ErrFmt(err)
	}
	pages, buttons, err := expandPermissionRows(rows)
	if err != nil {
		return nil, nil, stores.ErrFmt(err)
	}
	return pages, buttons, nil
}

func collectRoleIDs(roles []*authdto.RoleDto) []string {
	roleIDs := make([]string, 0, len(roles))
	for _, role := range roles {
		if role == nil || strings.TrimSpace(role.RoleID) == "" {
			continue
		}
		roleIDs = append(roleIDs, strings.TrimSpace(role.RoleID))
	}
	return roleIDs
}

func appendDefaultResources(resources []*authdto.ResourceDto) []*authdto.ResourceDto {
	uriMap := make(map[string]*authdto.ResourceDto, len(resources)+len(enums.DefaultAllowURIs))
	for _, resource := range resources {
		if resource == nil || strings.TrimSpace(resource.URI) == "" {
			continue
		}
		uriMap[resource.URI] = resource
	}
	appendDerivedResources(uriMap)
	for _, uri := range enums.DefaultAllowURIs {
		if _, ok := uriMap[uri]; ok {
			continue
		}
		uriMap[uri] = &authdto.ResourceDto{
			URI:     uri,
			Methods: defaultMethods(),
		}
	}
	result := make([]*authdto.ResourceDto, 0, len(uriMap))
	for _, resource := range uriMap {
		result = append(result, resource)
	}
	slices.SortFunc(result, func(a, b *authdto.ResourceDto) int {
		return strings.Compare(a.URI, b.URI)
	})
	return result
}

func appendDerivedResources(uriMap map[string]*authdto.ResourceDto) {
	routingManagement := uriMap["/routing-management"]
	if routingManagement == nil {
		return
	}
	if _, ok := uriMap["/kong-admin"]; ok {
		return
	}
	methods := append([]string(nil), routingManagement.Methods...)
	if len(methods) == 0 {
		methods = defaultMethods()
	}
	uriMap["/kong-admin"] = &authdto.ResourceDto{
		URI:     "/kong-admin",
		Methods: methods,
	}
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
