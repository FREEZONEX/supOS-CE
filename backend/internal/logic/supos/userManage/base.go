package userManage

import (
	"context"
	"sort"
	"strings"
	"time"

	"backend/internal/common/I18nUtils"
	cache "backend/internal/common/cache"
	"backend/internal/common/enums"
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/vo"
	authlogic "backend/internal/logic/supos/auth"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type baseUserManageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func newBaseUserManageLogic(ctx context.Context, svcCtx *svc.ServiceContext) baseUserManageLogic {
	return baseUserManageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l baseUserManageLogic) db() (*gorm.DB, error) {
	db := stores.GetCommonConn(l.ctx)
	if db == nil {
		return nil, errors.System.WithMsg("common database connection not initialized")
	}
	return db.WithContext(l.ctx), nil
}

func (l baseUserManageLogic) currentUser() *vo.UserInfoVo {
	if user, ok := l.ctx.Value(apiutil.UserKey).(*vo.UserInfoVo); ok {
		return user
	}
	return nil
}

func (l baseUserManageLogic) success() *types.OperationResult {
	return &types.OperationResult{Success: true}
}

func (l baseUserManageLogic) operatorID() string {
	if user := l.currentUser(); user != nil {
		return strings.TrimSpace(user.Sub)
	}
	return ""
}

func (l baseUserManageLogic) invalidateUserCache(userID string) {
	if userID == "" || cache.UserInfoCache == nil {
		return
	}
	cache.UserInfoCache.Delete(userID)
}

func (l baseUserManageLogic) updateUserCache(user *vo.UserInfoVo) {
	if user == nil || cache.UserInfoCache == nil {
		return
	}
	cache.UserInfoCache.Set(user.Sub, user)
}

func (l baseUserManageLogic) getIAMUser(db *gorm.DB, userID string) (*relationDB.IamUser, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.Parameter.WithMsg("userId is required")
	}
	var user relationDB.IamUser
	err := db.Where("id = ?", strings.TrimSpace(userID)).Take(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		l.Errorf("load iam user failed: %v", err)
		return nil, errors.System.WithMsg("failed to load user")
	}
	return &user, nil
}

func (l baseUserManageLogic) getIAMUserByUsername(db *gorm.DB, username string) (*relationDB.IamUser, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, nil
	}
	var user relationDB.IamUser
	err := db.Where("LOWER(username) = LOWER(?)", username).Take(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		l.Errorf("query user by username failed: %v", err)
		return nil, errors.System.WithMsg("failed to query user")
	}
	return &user, nil
}

func (l baseUserManageLogic) getIAMUserByEmail(db *gorm.DB, email string) (*relationDB.IamUser, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, nil
	}
	var user relationDB.IamUser
	err := db.Where("LOWER(email) = LOWER(?)", email).Take(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		l.Errorf("query user by email failed: %v", err)
		return nil, errors.System.WithMsg("failed to query user")
	}
	return &user, nil
}

func (l baseUserManageLogic) getRoleByID(db *gorm.DB, roleID string) (*relationDB.IamRole, error) {
	if strings.TrimSpace(roleID) == "" {
		return nil, errors.Parameter.WithMsg("role.id.empty")
	}
	var role relationDB.IamRole
	err := db.Where("id = ?", strings.TrimSpace(roleID)).Take(&role).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		l.Errorf("load role by id failed: %v", err)
		return nil, errors.System.WithMsg("failed to load role")
	}
	return &role, nil
}

func (l baseUserManageLogic) getRoleByName(db *gorm.DB, roleName string) (*relationDB.IamRole, error) {
	roleName = strings.TrimSpace(roleName)
	if roleName == "" {
		return nil, nil
	}
	if builtinKey := normalizeBuiltinRoleKey(roleName); builtinKey != "" {
		var role relationDB.IamRole
		err := db.Where("LOWER(role_key) = LOWER(?)", builtinKey).Take(&role).Error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		if err != nil {
			l.Errorf("load builtin role failed: %v", err)
			return nil, errors.System.WithMsg("failed to load role")
		}
		return &role, nil
	}
	var role relationDB.IamRole
	err := db.Where("LOWER(role_name) = LOWER(?) OR LOWER(role_key) = LOWER(?)", roleName, roleName).Take(&role).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		l.Errorf("load role by name failed: %v", err)
		return nil, errors.System.WithMsg("failed to load role")
	}
	return &role, nil
}

func (l baseUserManageLogic) listRoles(db *gorm.DB) ([]relationDB.IamRole, error) {
	var roles []relationDB.IamRole
	err := db.Where("status = 1").Order("builtin DESC, role_name ASC").Find(&roles).Error
	if err != nil {
		l.Errorf("load role list failed: %v", err)
		return nil, errors.System.WithMsg("failed to load roles")
	}
	return roles, nil
}

func (l baseUserManageLogic) countActiveRoles(db *gorm.DB) (int64, error) {
	var total int64
	if err := db.Model(&relationDB.IamRole{}).Where("status = 1").Count(&total).Error; err != nil {
		l.Errorf("count roles failed: %v", err)
		return 0, errors.System.WithMsg("failed to count roles")
	}
	return total, nil
}

func (l baseUserManageLogic) loadRolesForUsers(db *gorm.DB, userIDs []string) (map[string][]types.RoleSummary, error) {
	if len(userIDs) == 0 {
		return map[string][]types.RoleSummary{}, nil
	}
	type roleJoinRow struct {
		UserID      string `gorm:"column:user_id"`
		RoleID      string `gorm:"column:role_id"`
		RoleKey     string `gorm:"column:role_key"`
		RoleName    string `gorm:"column:role_name"`
		Description string `gorm:"column:description"`
	}
	var rows []roleJoinRow
	if err := db.Table("supos_user_role AS ur").
		Select("ur.user_id, r.id AS role_id, r.role_key, r.role_name, r.description").
		Joins("JOIN supos_role AS r ON r.id = ur.role_id").
		Where("ur.user_id IN ?", userIDs).
		Where("r.status = 1").
		Order("r.builtin DESC, r.role_name ASC").
		Scan(&rows).Error; err != nil {
		l.Errorf("load user roles failed: %v", err)
		return nil, errors.System.WithMsg("failed to load user roles")
	}
	result := make(map[string][]types.RoleSummary)
	for _, row := range rows {
		if row.UserID == "" || row.RoleID == "" {
			continue
		}
		display, desc := normalizeRoleDisplay(l.ctx, row.RoleKey, row.RoleName, row.Description)
		summary := types.RoleSummary{
			RoleID:          row.RoleID,
			RoleName:        firstNonEmpty(display, row.RoleName, row.RoleKey),
			RoleDescription: strings.TrimSpace(desc),
			ClientRole:      true,
		}
		result[row.UserID] = append(result[row.UserID], summary)
	}
	return result, nil
}

func (l baseUserManageLogic) getRoleResourceList(db *gorm.DB, roleID string) ([]types.RoleResource, error) {
	allow, _, err := l.getRolePermissionLists(db, roleID)
	return allow, err
}

func (l baseUserManageLogic) getRolePermissionLists(db *gorm.DB, roleID string) ([]types.RoleResource, []types.RoleResource, error) {
	allRows, err := l.loadAssignableResources(db)
	if err != nil {
		return nil, nil, err
	}

	assignedRows, err := l.loadRoleResources(db, roleID)
	if err != nil {
		return nil, nil, err
	}

	allow, deny := buildRolePermissionLists(allRows, assignedRows)
	return allow, deny, nil
}

func (l baseUserManageLogic) loadAssignableResources(db *gorm.DB) ([]relationDB.SuposResource, error) {
	var rows []relationDB.SuposResource
	err := db.Table("supos_resource AS sr").
		Select("sr.*").
		Where("sr.type IN ?", []int{2, 3, 5}).
		Where("COALESCE(sr.enable, true) = true").
		Order("sr.sort ASC, sr.id ASC").
		Scan(&rows).Error
	if err != nil {
		l.Errorf("load assignable resources failed: %v", err)
		return nil, errors.System.WithMsg("failed to load resources")
	}
	return rows, nil
}

func (l baseUserManageLogic) loadRoleResources(db *gorm.DB, roleID string) ([]relationDB.SuposResource, error) {
	var rows []relationDB.SuposResource
	err := db.Table("supos_resource AS sr").
		Select("sr.*").
		Joins("JOIN supos_role_resource AS rr ON rr.resource_id = sr.id").
		Where("rr.role_id = ?", strings.TrimSpace(roleID)).
		Where("COALESCE(sr.enable, true) = true").
		Order("sr.sort ASC, sr.id ASC").
		Scan(&rows).Error
	if err != nil {
		l.Errorf("load role resources failed: %v", err)
		return nil, errors.System.WithMsg("failed to load role resources")
	}
	return rows, nil
}

func (l baseUserManageLogic) mapRoleResourcesToIDs(
	db *gorm.DB,
	allowResources []types.RoleResource,
	denyResources []types.RoleResource,
) ([]int64, error) {
	rows, err := l.loadAssignableResources(db)
	if err != nil {
		return nil, err
	}
	return mapRoleResourcesToIDsFromRows(rows, allowResources, denyResources), nil
}

func (l baseUserManageLogic) replaceRoleResources(
	tx *gorm.DB,
	roleID string,
	allowResources []types.RoleResource,
	denyResources []types.RoleResource,
) error {
	resourceIDs, err := l.mapRoleResourcesToIDs(tx, allowResources, denyResources)
	if err != nil {
		return err
	}
	if err := tx.Where("role_id = ?", strings.TrimSpace(roleID)).Delete(&relationDB.IamRoleResource{}).Error; err != nil {
		l.Errorf("delete old role resources failed: %v", err)
		return errors.System.WithMsg("failed to save role resources")
	}
	if len(resourceIDs) == 0 {
		return nil
	}
	now := time.Now()
	items := make([]relationDB.IamRoleResource, 0, len(resourceIDs))
	operatorID := l.operatorID()
	for _, resourceID := range resourceIDs {
		items = append(items, relationDB.IamRoleResource{
			RoleID:     strings.TrimSpace(roleID),
			ResourceID: resourceID,
			CreatedAt:  now,
			CreatedBy:  operatorID,
		})
	}
	if err := tx.Create(&items).Error; err != nil {
		l.Errorf("create role resources failed: %v", err)
		return errors.System.WithMsg("failed to save role resources")
	}
	return nil
}

func (l baseUserManageLogic) resolveRoleIDs(db *gorm.DB, roles []types.RoleSummary) ([]string, error) {
	if len(roles) == 0 {
		return nil, nil
	}
	seen := make(map[string]struct{}, len(roles))
	result := make([]string, 0, len(roles))
	for _, role := range roles {
		var (
			roleModel *relationDB.IamRole
			err       error
		)
		if strings.TrimSpace(role.RoleID) != "" {
			roleModel, err = l.getRoleByID(db, role.RoleID)
		} else {
			roleModel, err = l.getRoleByName(db, role.RoleName)
		}
		if err != nil {
			return nil, err
		}
		if roleModel == nil || roleModel.Status != 1 {
			return nil, errors.Parameter.WithMsg("role.no.exist")
		}
		if _, ok := seen[roleModel.ID]; ok {
			continue
		}
		seen[roleModel.ID] = struct{}{}
		result = append(result, roleModel.ID)
	}
	return result, nil
}

func (l baseUserManageLogic) getDefaultUserRoleID(db *gorm.DB) (string, error) {
	role, err := l.getRoleByName(db, "user")
	if err != nil {
		return "", err
	}
	if role != nil {
		return role.ID, nil
	}
	return enums.RoleNormalUser.ID, nil
}

func (l baseUserManageLogic) replaceUserRoles(tx *gorm.DB, userID string, roleIDs []string) error {
	if err := tx.Where("user_id = ?", strings.TrimSpace(userID)).Delete(&relationDB.IamUserRole{}).Error; err != nil {
		l.Errorf("delete user roles failed: %v", err)
		return errors.System.WithMsg("failed to update user roles")
	}
	if len(roleIDs) == 0 {
		return nil
	}
	now := time.Now()
	operatorID := l.operatorID()
	items := make([]relationDB.IamUserRole, 0, len(roleIDs))
	for _, roleID := range dedupeStrings(roleIDs) {
		items = append(items, relationDB.IamUserRole{
			UserID:    strings.TrimSpace(userID),
			RoleID:    roleID,
			CreatedAt: now,
			CreatedBy: operatorID,
		})
	}
	if err := tx.Create(&items).Error; err != nil {
		l.Errorf("create user roles failed: %v", err)
		return errors.System.WithMsg("failed to update user roles")
	}
	return nil
}

func (l baseUserManageLogic) addUserRoles(tx *gorm.DB, userID string, roleIDs []string) error {
	if len(roleIDs) == 0 {
		return nil
	}
	now := time.Now()
	operatorID := l.operatorID()
	items := make([]relationDB.IamUserRole, 0, len(roleIDs))
	for _, roleID := range dedupeStrings(roleIDs) {
		items = append(items, relationDB.IamUserRole{
			UserID:    strings.TrimSpace(userID),
			RoleID:    roleID,
			CreatedAt: now,
			CreatedBy: operatorID,
		})
	}
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error; err != nil {
		l.Errorf("add user roles failed: %v", err)
		return errors.System.WithMsg("failed to update user roles")
	}
	return nil
}

func (l baseUserManageLogic) removeUserRoles(tx *gorm.DB, userID string, roleIDs []string) error {
	if len(roleIDs) == 0 {
		return nil
	}
	if err := tx.Where("user_id = ? AND role_id IN ?", strings.TrimSpace(userID), dedupeStrings(roleIDs)).
		Delete(&relationDB.IamUserRole{}).Error; err != nil {
		l.Errorf("remove user roles failed: %v", err)
		return errors.System.WithMsg("failed to update user roles")
	}
	return nil
}

func (l baseUserManageLogic) upsertUserPassword(tx *gorm.DB, userID, password string) error {
	hash, err := authlogic.HashPassword(strings.TrimSpace(password))
	if err != nil {
		l.Errorf("hash password failed: %v", err)
		return errors.System.WithMsg("failed to save password")
	}
	now := time.Now()
	if err := tx.Model(&relationDB.IamUser{}).
		Where("id = ?", strings.TrimSpace(userID)).
		Updates(map[string]any{
			"password":   hash,
			"updated_at": now,
		}).Error; err != nil {
		l.Errorf("save user password failed: %v", err)
		return errors.System.WithMsg("failed to save password")
	}
	return nil
}

func (l baseUserManageLogic) revokeUserSessions(tx *gorm.DB, userID string) error {
	now := time.Now()
	if err := tx.Model(&relationDB.IamSession{}).
		Where("user_id = ? AND revoked_at IS NULL", strings.TrimSpace(userID)).
		Update("revoked_at", now).Error; err != nil {
		l.Errorf("revoke user sessions failed: %v", err)
		return errors.System.WithMsg("failed to revoke user sessions")
	}
	return nil
}

func (l baseUserManageLogic) protectedRole(role *relationDB.IamRole) bool {
	if role == nil {
		return false
	}
	if role.Builtin {
		return true
	}
	switch normalizeBuiltinRoleKey(firstNonEmpty(role.RoleKey, role.RoleName)) {
	case "admin", "user":
		return true
	default:
		return false
	}
}

func normalizeBuiltinRoleKey(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "admin", "super-admin", strings.ToLower(enums.RoleAdmin.Comment), strings.ToLower(enums.RoleSuperAdmin.Comment):
		return "admin"
	case "user", "normal-user", strings.ToLower(enums.RoleNormalUser.Comment):
		return "user"
	default:
		return ""
	}
}

func normalizeRoleKey(name string) string {
	if builtin := normalizeBuiltinRoleKey(name); builtin != "" {
		return builtin
	}
	return strings.ToLower(strings.TrimSpace(name))
}

func normalizeRoleDisplay(ctx context.Context, roleKey, roleName, description string) (display, desc string) {
	switch normalizeBuiltinRoleKey(firstNonEmpty(roleKey, roleName)) {
	case "admin":
		return I18nUtils.GetMessageWithCtx(ctx, enums.RoleAdmin.I18nCode), description
	case "user":
		return I18nUtils.GetMessageWithCtx(ctx, enums.RoleNormalUser.I18nCode), description
	default:
		if strings.TrimSpace(roleName) != "" {
			return strings.TrimSpace(roleName), description
		}
		return strings.TrimSpace(roleKey), description
	}
}

func buildPagePermission(resource relationDB.SuposResource) string {
	if resource.URLType != nil && *resource.URLType == 1 && resource.URL != nil && strings.TrimSpace(*resource.URL) != "" {
		return strings.TrimSpace(*resource.URL)
	}
	if strings.TrimSpace(resource.Code) == "" {
		return ""
	}
	return "/" + strings.TrimSpace(resource.Code)
}

func normalizePermissionURI(uri string) string {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return ""
	}
	if idx := strings.Index(uri, "$"); idx >= 0 {
		return strings.TrimSpace(uri[:idx])
	}
	return uri
}

func defaultMethods() []string {
	return []string{"get", "post", "put", "delete", "patch", "head", "options"}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func dedupeStrings(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		key := strings.TrimSpace(item)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func boolValue(val *bool, def bool) bool {
	if val == nil {
		return def
	}
	return *val
}

func intValue(val *int, def int) int {
	if val == nil {
		return def
	}
	return *val
}
