package authsvc

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	cache "backend/internal/common/cache"
	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/vo"
	keycloakrepo "backend/internal/repo/keycloak"
	"backend/share/clients"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

// FetchUserInfo returns the user profile associated with the access token.
// The result may be served from cache when allowCache is true.
func FetchUserInfo(ctx context.Context, kc *clients.KeycloakClient, accessToken string, allowCache bool, defaultHome, realm string) (*vo.UserInfoVo, string, error) {
	if kc == nil {
		return nil, "", errors.System.WithMsg("keycloak client not configured")
	}
	if strings.TrimSpace(accessToken) == "" {
		return nil, "", errors.Parameter.WithMsg("access token is empty")
	}

	claims, err := DecodeJWTClaims(accessToken)
	if err != nil {
		return nil, "", err
	}
	sub := ClaimString(claims, "sub")
	if sub == "" {
		return nil, "", errors.Parameter.WithMsg("token payload missing sub")
	}

	if allowCache && cache.UserInfoCache != nil {
		if cached, ok := cache.UserInfoCache.Get(sub); ok && cached != nil {
			return cached, sub, nil
		}
	}

	if keycloakrepo.Enabled() {
		repo, err := keycloakrepo.NewAuthRepo(ctx)
		if err != nil {
			logx.WithContext(ctx).Errorf("init keycloak repo failed: %v", err)
		} else if repo != nil {
			if user, err := repo.BuildUserInfo(ctx, realm, sub, defaultHome); err == nil && user != nil {
				if err := updateFirstLogin(ctx, kc, user); err != nil {
					logx.WithContext(ctx).Errorf("update first login attributes failed: %v", err)
				}
				if cache.UserInfoCache != nil {
					cache.UserInfoCache.Set(sub, user)
				}
				return user, sub, nil
			} else if err != nil {
				logx.WithContext(ctx).Errorf("load user info from database failed, fallback to keycloak api: %v", err)
			}
		}
	}
	logx.WithContext(ctx).Info("keycloakrepo not Enabled")
	user, err := loadUserInfoFromKeycloak(ctx, kc, sub, claims, accessToken, defaultHome)
	if err != nil {
		return nil, sub, err
	}
	if err := updateFirstLogin(ctx, kc, user); err != nil {
		logx.WithContext(ctx).Errorf("update first login attributes failed: %v", err)
	}
	if user != nil && cache.UserInfoCache != nil {
		cache.UserInfoCache.Set(sub, user)
	}
	return user, sub, nil
}

func loadUserInfoFromKeycloak(ctx context.Context, kc *clients.KeycloakClient, sub string, claims map[string]any, accessToken, defaultHome string) (*vo.UserInfoVo, error) {
	info, err := kc.UserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	preferredUsername := ClaimString(claims, "preferred_username")
	if preferredUsername == "" && info != nil {
		preferredUsername = info.Username
	}
	if preferredUsername == "" {
		preferredUsername = sub
	}

	user := vo.NewUserInfoVo(sub, preferredUsername)
	if info != nil {
		user.Email = trimString(info.Email)
		user.Enabled = info.Enabled
		user.FirstName = trimString(info.FirstName)
	}
	if user.HomePage == "" {
		user.HomePage = defaultHome
	}

	// Enrich with admin user profile if possible.
	adminProfile, err := kc.FetchUser(preferredUsername)
	if err == nil && adminProfile != nil {
		if user.Email == "" {
			user.Email = trimString(adminProfile.Email)
		}
		if adminProfile.Attributes != nil {
			applyUserAttributes(user, adminProfile.Attributes)
		}
	}

	roles, err := loadRoles(kc, sub)
	if err != nil {
		logx.WithContext(ctx).Errorf("load user roles failed: %v", err)
	} else {
		user.RoleList = roles
	}
	user.SuperAdmin = user.IsSuperAdmin()
	return user, nil
}

func loadRoles(kc *clients.KeycloakClient, userID string) ([]*authdto.RoleDto, error) {
	rawRoles, err := kc.GetRoleListByUserID(userID)
	if err != nil {
		return nil, err
	}
	if rawRoles == nil {
		return nil, nil
	}
	data, err := json.Marshal(rawRoles)
	if err != nil {
		return nil, err
	}
	var roleMappings struct {
		RealmMappings  []clients.KeycloakRoleInfoDto `json:"realmMappings"`
		ClientMappings map[string]struct {
			Mappings []clients.KeycloakRoleInfoDto `json:"mappings"`
		} `json:"clientMappings"`
	}
	if err := json.Unmarshal(data, &roleMappings); err != nil {
		return nil, err
	}

	var roles []*authdto.RoleDto
	for _, r := range roleMappings.RealmMappings {
		roles = append(roles, convertRole(r))
	}
	for _, cm := range roleMappings.ClientMappings {
		for _, r := range cm.Mappings {
			roles = append(roles, convertRole(r))
		}
	}
	return roles, nil
}

func convertRole(role clients.KeycloakRoleInfoDto) *authdto.RoleDto {
	return &authdto.RoleDto{
		RoleID:          role.ID,
		RoleName:        role.Name,
		RoleDescription: trimString(role.Description),
		ClientRole:      role.ClientRole,
	}
}

func applyUserAttributes(user *vo.UserInfoVo, attrs map[string]any) {
	if user == nil || attrs == nil {
		return
	}
	// if v := attributeString(attrs, "homePage"); v != "" {
	// 	user.HomePage = v
	// }
	if v := attributeString(attrs, "phone"); v != "" {
		user.Phone = v
	}
	if v := attributeString(attrs, "source"); v != "" {
		user.Source = v
	}
	if v := attributeString(attrs, "firstTimeLogin"); v != "" {
		if iv, err := strconv.Atoi(v); err == nil {
			user.FirstTimeLogin = iv
		}
	}
	if v := attributeString(attrs, "tipsEnable"); v != "" {
		if iv, err := strconv.Atoi(v); err == nil {
			user.TipsEnable = iv
		}
	}
}

func attributeString(attrs map[string]any, key string) string {
	raw, ok := attrs[key]
	if !ok || raw == nil {
		return ""
	}
	switch v := raw.(type) {
	case string:
		return strings.TrimSpace(v)
	case []string:
		if len(v) > 0 {
			return strings.TrimSpace(v[0])
		}
	case []any:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return strings.TrimSpace(s)
			}
		}
	}
	return ""
}

func trimString(val string) string {
	return strings.TrimSpace(val)
}

func updateFirstLogin(ctx context.Context, kc *clients.KeycloakClient, user *vo.UserInfoVo) error {
	if kc == nil || user == nil || strings.TrimSpace(user.Sub) == "" {
		return nil
	}
	if user.FirstTimeLogin != 1 {
		return nil
	}

	attrs := map[string]string{
		"firstTimeLogin": "0",
		"phone":          trimString(user.Phone),
		"tipsEnable":     "0",
		"source":         trimString(user.Source),
	}

	payload := map[string]any{
		"attributes": convertAttributesPayload(attrs),
	}
	if email := trimString(user.Email); email != "" {
		payload["email"] = email
	}

	return kc.UpdateUser(user.Sub, payload)
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
