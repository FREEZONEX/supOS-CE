package keycloak

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/enums"
	"backend/internal/common/utils/langutil"
	"backend/internal/common/vo"
	"backend/internal/repo/relationDB"

	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// AuthRepo provides read access to Keycloak tables.
type AuthRepo struct {
	db *gorm.DB
}

// NewAuthRepo creates a repository bound to the provided context.
func NewAuthRepo(ctx context.Context) (*AuthRepo, error) {
	conn := GetConn(ctx)
	if conn == nil {
		if isConfigured {
			return nil, stores.ErrFmt(fmt.Errorf("keycloak database connection not initialized"))
		}
		return nil, nil
	}
	return &AuthRepo{db: conn}, nil
}

type userEntity struct {
	ID            string `gorm:"column:id"`
	Username      string `gorm:"column:username"`
	Email         string `gorm:"column:email"`
	FirstName     string `gorm:"column:first_name"`
	LastName      string `gorm:"column:last_name"`
	Enabled       bool   `gorm:"column:enabled"`
	EmailVerified bool   `gorm:"column:email_verified"`
	RealmID       string `gorm:"column:realm_id"`
}

func (userEntity) TableName() string { return "user_entity" }

type attributeRow struct {
	Name  string `gorm:"column:name"`
	Value string `gorm:"column:value"`
}

type resourceRow struct {
	PolicyID   string `gorm:"column:policy_id"`
	ResourceID string `gorm:"column:resource_id"`
	URI        string `gorm:"column:uri"`
}

// BuildUserInfo assembles the user profile for the given user id.
func (r *AuthRepo) BuildUserInfo(ctx context.Context, realm string, userID string, defaultHome string) (*vo.UserInfoVo, error) {
	if r == nil || r.db == nil {
		return nil, stores.ErrFmt(fmt.Errorf("nil repository"))
	}

	var entity userEntity
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&entity).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}

	user := vo.NewUserInfoVo(entity.ID, entity.Username)
	user.Email = strings.TrimSpace(entity.Email)
	user.Enabled = entity.Enabled
	user.EmailVerified = entity.EmailVerified
	user.FirstName = strings.TrimSpace(entity.FirstName)
	if user.HomePage == "" {
		user.HomePage = defaultHome
	}

	attrs, err := r.getUserAttributes(ctx, userID)
	if err != nil {
		logx.WithContext(ctx).Errorf("BuildUserInfo getUserAttributes err %v", err)
		return nil, err
	}
	applyAttributesFromMap(user, attrs)

	roles, err := r.getRoles(ctx, realm, userID)
	if err != nil {
		logx.WithContext(ctx).Errorf("BuildUserInfo getRoles err %v", err)
		return nil, err
	}

	user.RoleList = filterVisibleRoles(roles)

	denyRoleIDs := collectRoleIDsByPrefix(roles, "deny")
	allowRoleIDs := collectAllowedRoleIDs(roles)

	denyResources, err := r.getResourcesByRoleIDs(ctx, denyRoleIDs)
	if err != nil {
		logx.WithContext(ctx).Errorf("load deny resources failed: %v", err)
		denyResources = nil
	}

	allowResources, err := r.getResourcesByRoleIDs(ctx, allowRoleIDs)
	if err != nil {
		logx.WithContext(ctx).Errorf("load allow resources failed: %v", err)
		allowResources = nil
	}

	user.DenyResourceList = sanitizeResources(denyResources)
	user.ResourceList = finalizeAllowResources(allowResources)

	lang, langErr := resolveUserMainLanguage(ctx, userID)
	if langErr != nil {
		logx.WithContext(ctx).Errorf("resolve user(%s) language failed: %v", userID, langErr)
	}
	user.MainLanguage = lang
	user.SuperAdmin = user.IsSuperAdmin()
	return user, nil
}

// GetRoleAllowResources returns sanitized allow resources for the given role.
func (r *AuthRepo) GetRoleAllowResources(ctx context.Context, roleID string) ([]*authdto.ResourceDto, error) {
	if roleID == "" {
		return nil, nil
	}
	resources, err := r.getResourcesByRoleIDs(ctx, []string{roleID})
	if err != nil {
		return nil, err
	}
	return finalizeAllowResources(resources), nil
}

// GetRoleDenyResources returns sanitized deny resources for the given role.
func (r *AuthRepo) GetRoleDenyResources(ctx context.Context, roleID string) ([]*authdto.ResourceDto, error) {
	if roleID == "" {
		return nil, nil
	}
	resources, err := r.getResourcesByRoleIDs(ctx, []string{roleID})
	if err != nil {
		return nil, err
	}
	return sanitizeResources(resources), nil
}

func (r *AuthRepo) getUserAttributes(ctx context.Context, userID string) (map[string]string, error) {
	var rows []attributeRow
	err := r.db.WithContext(ctx).Table("user_attribute").Where("user_id = ?", userID).Find(&rows).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		if row.Name == "" {
			continue
		}
		result[row.Name] = row.Value
	}
	return result, nil
}

func (r *AuthRepo) getRoles(ctx context.Context, realm, userID string) ([]*authdto.RoleDto, error) {
	if realm == "" {
		return nil, stores.ErrFmt(fmt.Errorf("realm cannot be empty"))
	}

	var baseRoles []*authdto.RoleDto
	roleQuery := `
SELECT r.id AS role_id,
       r.name AS role_name,
       r.description AS role_description,
       r.client_role
FROM user_entity ue
         JOIN user_role_mapping urm ON ue.id = urm.user_id
         JOIN keycloak_role r ON urm.role_id = r.id
         JOIN realm realm ON ue.realm_id = realm.id
WHERE ue.id = ? AND realm.name = ?
`
	if err := r.db.WithContext(ctx).Raw(roleQuery, userID, realm).Scan(&baseRoles).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}

	compositeIDs := collectCompositeRoleIDs(baseRoles)
	childRoles, err := r.getChildRoles(ctx, compositeIDs)
	if err != nil {
		return nil, err
	}

	allRoles := mergeRoles(baseRoles, childRoles)
	return allRoles, nil
}

func (r *AuthRepo) getChildRoles(ctx context.Context, compositeIDs []string) ([]*authdto.RoleDto, error) {
	if len(compositeIDs) == 0 {
		return nil, nil
	}
	query := `
SELECT r.id AS role_id,
       r.name AS role_name,
       r.description AS role_description,
       r.client_role
FROM keycloak_role r
WHERE r.client = (SELECT id FROM client WHERE client_id = 'tier0')
  AND r.id IN (
      SELECT child_role FROM composite_role WHERE composite IN ?
  )
`
	var roles []*authdto.RoleDto
	if err := r.db.WithContext(ctx).Raw(query, compositeIDs).Scan(&roles).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}
	return roles, nil
}

func (r *AuthRepo) getResourcesByRoleIDs(ctx context.Context, roleIDs []string) ([]*authdto.ResourceDto, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	policySet := make(map[string]struct{})
	for _, roleID := range roleIDs {
		ids, err := r.getPolicyIDsByRoleID(ctx, roleID)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			policySet[id] = struct{}{}
		}
	}
	if len(policySet) == 0 {
		return nil, nil
	}
	policyIDs := make([]string, 0, len(policySet))
	for id := range policySet {
		policyIDs = append(policyIDs, id)
	}
	return r.getResourceListByPolicyIDs(ctx, policyIDs)
}

func (r *AuthRepo) getPolicyIDsByRoleID(ctx context.Context, roleID string) ([]string, error) {
	sql := `
SELECT DISTINCT t.policy_id
FROM policy_config t
WHERE t.value IS NOT NULL
  AND t.value <> ''
  AND t.name = 'roles'
  AND t.value::jsonb IS NOT NULL
  AND t.value::jsonb @> ?::jsonb
`
	payload := fmt.Sprintf(`[{"id":"%s"}]`, roleID)
	var ids []string
	if err := r.db.WithContext(ctx).Raw(sql, payload).Scan(&ids).Error; err != nil {
		return nil, stores.ErrFmt(err)
	}
	return ids, nil
}

func (r *AuthRepo) getResourceListByPolicyIDs(ctx context.Context, policyIDs []string) ([]*authdto.ResourceDto, error) {
	if len(policyIDs) == 0 {
		return nil, nil
	}
	var rows []resourceRow
	err := r.db.WithContext(ctx).
		Table("associated_policy AS ap").
		Select(`rp.policy_id, rp.resource_id, ru."value" AS uri`).
		Joins("LEFT JOIN resource_policy rp ON ap.policy_id = rp.policy_id").
		Joins(`LEFT JOIN resource_uris ru ON ru.resource_id = rp.resource_id`).
		Where(`ru."value" <> ''`).
		Where("ap.associated_policy_id IN ?", policyIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	resources := make([]*authdto.ResourceDto, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.URI) == "" {
			continue
		}
		resources = append(resources, &authdto.ResourceDto{
			PolicyID:   row.PolicyID,
			ResourceID: row.ResourceID,
			URI:        row.URI,
		})
	}
	return resources, nil
}

func applyAttributesFromMap(user *vo.UserInfoVo, attrs map[string]string) {
	if user == nil || len(attrs) == 0 {
		return
	}
	if v := strings.TrimSpace(attrs["homePage"]); v != "" {
		user.HomePage = v
	}
	if v := strings.TrimSpace(attrs["phone"]); v != "" {
		user.Phone = v
	}
	if v := strings.TrimSpace(attrs["source"]); v != "" {
		user.Source = v
	}
	if v := strings.TrimSpace(attrs["firstTimeLogin"]); v != "" {
		if iv, err := strconv.Atoi(v); err == nil {
			user.FirstTimeLogin = iv
		} else {
			user.FirstTimeLogin = 1
		}
	}
	user.TipsEnable = user.FirstTimeLogin
	// if v := strings.TrimSpace(attrs["tipsEnable"]); v != "" {
	// 	if iv, err := strconv.Atoi(v); err == nil {
	// 		user.TipsEnable = iv
	// 	}
	// }
}

func collectCompositeRoleIDs(roles []*authdto.RoleDto) []string {
	var ids []string
	for _, role := range roles {
		if role != nil && !role.ClientRole {
			ids = append(ids, role.RoleID)
		}
	}
	return ids
}

func mergeRoles(base, extra []*authdto.RoleDto) []*authdto.RoleDto {
	result := make([]*authdto.RoleDto, 0, len(base)+len(extra))
	result = append(result, base...)
	result = append(result, extra...)
	unique := make(map[string]*authdto.RoleDto)
	for _, role := range result {
		if role == nil {
			continue
		}
		unique[role.RoleID] = role
	}
	out := make([]*authdto.RoleDto, 0, len(unique))
	for _, role := range unique {
		out = append(out, role)
	}
	return out
}

func filterVisibleRoles(roles []*authdto.RoleDto) []*authdto.RoleDto {
	var filtered []*authdto.RoleDto
	for _, role := range roles {
		if role == nil {
			continue
		}
		if !role.ClientRole || enums.IsIgnoredRoleID(role.RoleID) || enums.IsIgnoredRoleName(role.RoleName) || strings.HasPrefix(role.RoleName, "deny-") {
			continue
		}
		filtered = append(filtered, role)
	}
	return filtered
}

func collectRoleIDsByPrefix(roles []*authdto.RoleDto, prefix string) []string {
	var ids []string
	for _, role := range roles {
		if role == nil {
			continue
		}
		if strings.HasPrefix(role.RoleName, prefix) {
			ids = append(ids, role.RoleID)
		}
	}
	return ids
}

func collectAllowedRoleIDs(roles []*authdto.RoleDto) []string {
	var ids []string
	for _, role := range roles {
		if role == nil {
			continue
		}
		if strings.HasPrefix(role.RoleName, "deny") {
			continue
		}
		if enums.IsIgnoredRoleName(role.RoleName) {
			continue
		}
		ids = append(ids, role.RoleID)
	}
	return ids
}

func sanitizeResources(resources []*authdto.ResourceDto) []*authdto.ResourceDto {
	if len(resources) == 0 {
		return nil
	}
	for _, res := range resources {
		if res == nil {
			continue
		}
		res.URI = removeIfUriSuffix(res.URI)
		res.Methods = transMethodList(res.URI)
	}
	return deduplicateResources(resources)
}

func finalizeAllowResources(resources []*authdto.ResourceDto) []*authdto.ResourceDto {
	return appendCommonResources(sanitizeResources(resources))
}

func appendCommonResources(resources []*authdto.ResourceDto) []*authdto.ResourceDto {
	uriMap := make(map[string]*authdto.ResourceDto)
	for _, res := range resources {
		if res == nil {
			continue
		}
		uriMap[res.URI] = res
	}
	for _, uri := range enums.DefaultAllowURIs {
		if _, exists := uriMap[uri]; exists {
			continue
		}
		uriMap[uri] = &authdto.ResourceDto{
			URI:     uri,
			Methods: transMethodList(uri),
		}
	}
	out := make([]*authdto.ResourceDto, 0, len(uriMap))
	for _, res := range uriMap {
		out = append(out, res)
	}
	return out
}

var (
	defMethods = []string{"get", "post", "put", "delete", "patch", "head", "options"}
)

func transMethodList(uri string) []string {
	if !strings.Contains(uri, "$") {
		return defMethods
	}
	methodStr := uri[strings.Index(uri, "$")+1:]
	if strings.TrimSpace(methodStr) == "" {
		return defMethods
	}
	items := strings.Split(methodStr, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(strings.ToLower(item))
		if item != "" {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return defMethods
	}
	return result
}

func removeIfUriSuffix(uri string) string {
	if idx := strings.Index(uri, "$"); idx != -1 {
		return uri[:idx]
	}
	return uri
}

func deduplicateResources(resources []*authdto.ResourceDto) []*authdto.ResourceDto {
	if len(resources) == 0 {
		return nil
	}
	result := make(map[string]*authdto.ResourceDto, len(resources))
	for _, res := range resources {
		if res == nil {
			continue
		}
		key := res.URI
		if existing, ok := result[key]; ok {
			if len(existing.Methods) >= len(res.Methods) {
				continue
			}
		}
		result[key] = res
	}
	out := make([]*authdto.ResourceDto, 0, len(result))
	for _, res := range result {
		out = append(out, res)
	}
	return out
}

func resolveUserMainLanguage(ctx context.Context, userID string) (string, error) {
	defaultLang := langutil.SystemLocale()
	if strings.TrimSpace(userID) == "" {
		return defaultLang, nil
	}

	lang, err := loadUserMainLanguage(ctx, userID)
	if err != nil {
		return defaultLang, err
	}
	if lang == "" {
		return defaultLang, nil
	}
	return lang, nil
}

func loadUserMainLanguage(ctx context.Context, userID string) (string, error) {
	conn := stores.GetCommonConn(ctx)
	if conn == nil {
		return "", stores.ErrFmt(fmt.Errorf("common database connection not initialized"))
	}

	var cfg relationDB.UnsPersonConfig
	err := conn.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("update_at DESC").
		First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", stores.ErrFmt(err)
	}
	return strings.TrimSpace(cfg.MainLanguage), nil
}
