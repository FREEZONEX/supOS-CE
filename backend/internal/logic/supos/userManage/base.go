package userManage

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"backend/internal/common/I18nUtils"
	cache "backend/internal/common/cache"
	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/enums"
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/vo"
	keycloakrepo "backend/internal/repo/keycloak"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/clients"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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

func (l baseUserManageLogic) keycloakClient() (*clients.KeycloakClient, error) {
	if l.svcCtx == nil || l.svcCtx.Keycloak == nil {
		return nil, errors.System.WithMsg("keycloak client not configured")
	}
	return l.svcCtx.Keycloak, nil
}

func (l baseUserManageLogic) keycloakDB() (*gorm.DB, error) {
	db := keycloakrepo.GetConn(l.ctx)
	if db == nil {
		if keycloakrepo.Enabled() {
			return nil, errors.System.WithMsg("keycloak database connection not initialized")
		}
		return nil, errors.System.WithMsg("keycloak database not configured")
	}
	return db, nil
}

func (l baseUserManageLogic) authRepo() (*keycloakrepo.AuthRepo, error) {
	return keycloakrepo.NewAuthRepo(l.ctx)
}

func (l baseUserManageLogic) currentUser() *vo.UserInfoVo {
	if user, ok := l.ctx.Value(apiutil.UserKey).(*vo.UserInfoVo); ok {
		return user
	}
	return nil
}

func (l baseUserManageLogic) realm() string {
	if l.svcCtx == nil {
		return ""
	}
	return l.svcCtx.Config.OAuthKeyCloak.Realm
}

func (l baseUserManageLogic) success() *types.OperationResult {
	return &types.OperationResult{Success: true}
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

type kcUserEntity struct {
	ID        string `gorm:"column:id"`
	Username  string `gorm:"column:username"`
	FirstName string `gorm:"column:first_name"`
	Email     string `gorm:"column:email"`
	Enabled   bool   `gorm:"column:enabled"`
}

type kcAttributeRow struct {
	UserID string `gorm:"column:user_id"`
	Name   string `gorm:"column:name"`
	Value  string `gorm:"column:value"`
}

type kcRoleRow struct {
	UserID      string `gorm:"column:user_id"`
	RoleID      string `gorm:"column:role_id"`
	RoleName    string `gorm:"column:role_name"`
	Description string `gorm:"column:role_description"`
	ClientRole  bool   `gorm:"column:client_role"`
}

func (l baseUserManageLogic) getUserEntity(userID string) (*kcUserEntity, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.Parameter.WithMsg("userId is empty")
	}
	db, err := l.keycloakDB()
	if err != nil {
		return nil, err
	}
	var entity kcUserEntity
	if err := db.Table("user_entity").
		Select("id, username, first_name, email, enabled").
		Where("id = ?", userID).
		Take(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		l.Errorf("load user entity failed: %v", err)
		return nil, errors.System.WithMsg("failed to load user info")
	}
	return &entity, nil
}

func (l baseUserManageLogic) loadAttributesForUsers(db *gorm.DB, userIDs []string) (map[string]map[string]string, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var rows []kcAttributeRow
	if err := db.Table("user_attribute").
		Select("user_id, name, value").
		Where("user_id IN ?", userIDs).
		Find(&rows).Error; err != nil {
		l.Errorf("load user attributes failed: %v", err)
		return nil, errors.System.WithMsg("failed to load user attributes")
	}
	result := make(map[string]map[string]string)
	for _, row := range rows {
		if row.UserID == "" || row.Name == "" {
			continue
		}
		m, ok := result[row.UserID]
		if !ok {
			m = make(map[string]string)
			result[row.UserID] = m
		}
		m[row.Name] = row.Value
	}
	return result, nil
}

func (l baseUserManageLogic) loadRolesForUsers(db *gorm.DB, userIDs []string) (map[string][]types.RoleSummary, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var rows []kcRoleRow
	if err := db.Table("user_role_mapping AS urm").
		Select("urm.user_id, r.id AS role_id, r.name AS role_name, r.description AS role_description, r.client_role").
		Joins("JOIN keycloak_role AS r ON r.id = urm.role_id").
		Where("urm.user_id IN ?", userIDs).
		Find(&rows).Error; err != nil {
		l.Errorf("load user roles failed: %v", err)
		return nil, errors.System.WithMsg("failed to load user roles")
	}
	result := make(map[string][]types.RoleSummary)
	for _, row := range rows {
		if row.UserID == "" || row.RoleID == "" {
			continue
		}
		if !row.ClientRole || enums.IsIgnoredRoleID(row.RoleID) || enums.IsIgnoredRoleName(row.RoleName) || strings.HasPrefix(row.RoleName, "deny-") {
			continue
		}
		display, desc := normalizeRoleDisplay(l.ctx, row.RoleID, row.RoleName, row.Description)
		summary := types.RoleSummary{
			RoleID:          row.RoleID,
			RoleName:        strings.TrimSpace(display),
			RoleDescription: strings.TrimSpace(desc),
			ClientRole:      row.ClientRole,
		}
		if summary.RoleName == "" {
			summary.RoleName = strings.TrimSpace(row.RoleName)
		}
		result[row.UserID] = append(result[row.UserID], summary)
	}
	return result, nil
}

func convertAttributesPayload(attrs map[string]string) map[string]any {
	if len(attrs) == 0 {
		return map[string]any{}
	}
	payload := make(map[string]any, len(attrs))
	for key, val := range attrs {
		if strings.TrimSpace(key) == "" {
			continue
		}
		payload[key] = []string{val}
	}
	return payload
}

func toRoleResourceList(resources []*authdto.ResourceDto) []types.RoleResource {
	if len(resources) == 0 {
		return nil
	}
	result := make([]types.RoleResource, 0, len(resources))
	for _, res := range resources {
		if res == nil {
			continue
		}
		result = append(result, types.RoleResource{
			PolicyID:   res.PolicyID,
			ResourceID: res.ResourceID,
			URI:        res.URI,
			Methods:    slices.Clone(res.Methods),
		})
	}
	return result
}

func normalizeRoleDisplay(ctx context.Context, roleID, roleName, description string) (display, desc string) {
	if enumRole, ok := enums.RoleParse(roleID); ok {
		return enumsDisplay(ctx, enumRole), description
	}
	if strings.TrimSpace(description) != "" {
		return description, description
	}
	return roleName, description
}

func enumsDisplay(ctx context.Context, role enums.RoleEnum) string {
	return I18nUtils.GetMessageWithCtx(ctx, role.I18nCode)
}

func stringSlice(values []types.RoleResource) []string {
	if len(values) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		uri := strings.TrimSpace(v.URI)
		if uri == "" {
			continue
		}
		set[uri] = struct{}{}
	}
	result := make([]string, 0, len(set))
	for uri := range set {
		result = append(result, uri)
	}
	slices.Sort(result)
	return result
}

func dedupeStrings(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		key := strings.TrimSpace(item)
		if key == "" {
			continue
		}
		set[key] = struct{}{}
	}
	result := make([]string, 0, len(set))
	for item := range set {
		result = append(result, item)
	}
	slices.Sort(result)
	return result
}

func boolValue(val *bool) bool {
	if val == nil {
		return false
	}
	return *val
}

func intValue(val *int, def int) int {
	if val == nil {
		return def
	}
	return *val
}

type RolesReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (l baseUserManageLogic) applyUserRoles(userID string, roles []types.RoleSummary, add bool) error {
	if userID == "" || len(roles) == 0 {
		return nil
	}
	kc, err := l.keycloakClient()
	if err != nil {
		return err
	}
	roleMap := make(map[string]clients.KeycloakRoleInfoDto)
	for _, summary := range roles {
		roleInfo, err := l.resolveRole(kc, summary)
		if err != nil {
			return err
		}
		if roleInfo == nil {
			continue
		}
		roleMap[roleInfo.ID] = *roleInfo
		if roleInfo.ID != enums.RoleSuperAdmin.ID {
			denyName := fmt.Sprintf("deny-%s", roleInfo.Name)
			denyRole, err := kc.FetchRole(denyName)
			if err != nil {
				l.Errorf("failed to fetch deny role %s: %v", denyName, err)
				return errors.System.WithMsg("failed to load deny role")
			}
			if denyRole != nil {
				roleMap[denyRole.ID] = *denyRole
			}
		}
	}
	if len(roleMap) == 0 {
		return nil
	}
	values := make([]clients.KeycloakRoleInfoDto, 0, len(roleMap))
	for _, info := range roleMap {
		values = append(values, info)
	}
	action := "assign"
	if !add {
		action = "remove"
	}
	if err := kc.SetUserRoles(userID, values, add); err != nil {
		l.Errorf("failed to %s roles for user %s: %v", action, userID, err)
		return errors.System.WithMsg("failed to update user roles")
	}
	return nil
}

func (l baseUserManageLogic) resolveRole(kc *clients.KeycloakClient, summary types.RoleSummary) (*clients.KeycloakRoleInfoDto, error) {
	if kc == nil {
		return nil, errors.System.WithMsg("keycloak client not configured")
	}
	if id := strings.TrimSpace(summary.RoleID); id != "" {
		role, err := kc.FetchRoleByID(id)
		if err != nil {
			l.Errorf("failed to fetch role by id %s: %v", id, err)
			return nil, errors.System.WithMsg("failed to load role")
		}
		return role, nil
	}
	name := strings.TrimSpace(summary.RoleName)
	if name == "" {
		return nil, nil
	}
	role, err := kc.FetchRole(name)
	if err != nil {
		l.Errorf("failed to fetch role by name %s: %v", name, err)
		return nil, errors.System.WithMsg("failed to load role")
	}
	if role == nil && (name == enums.RoleSuperAdmin.Comment || name == enums.RoleSuperAdmin.Name) {
		role, err = kc.FetchRoleByID(enums.RoleSuperAdmin.ID)
		if err != nil {
			l.Errorf("failed to fetch super admin role: %v", err)
			return nil, errors.System.WithMsg("failed to load role")
		}
	}
	return role, nil
}
