// Code scaffolded by goctl. Safe to edit.

package openservice

import (
	"context"
	"strconv"

	authdto "backend/internal/common/dto/auth"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRoleService creates a new RoleService instance
func NewRoleService(ctx context.Context, svcCtx *svc.ServiceContext) *RoleService {
	return &RoleService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetRoleListByUserId 获取用户的角色列表
func (s *RoleService) GetRoleListByUserId(userID string) ([]*authdto.RoleDto, error) {
	if userID == "" {
		return nil, nil
	}

	if s.svcCtx.Keycloak == nil {
		// Keycloak 未配置，返回空列表
		return []*authdto.RoleDto{}, nil
	}

	// 通过 KeycloakClient 获取用户角色映射
	roleMappings, err := s.svcCtx.Keycloak.GetRoleListByUserID(userID)
	if err != nil {
		s.Errorf("Failed to get role mappings for user %s: %v", userID, err)
		return nil, err
	}

	// 解析 clientMappings 并转换为 RoleDto 列表
	var roles []*authdto.RoleDto
	if clientMappings, ok := roleMappings["clientMappings"].(map[string]any); ok {
		for _, clientRoles := range clientMappings {
			if rolesList, ok := clientRoles.([]any); ok {
				for _, roleItem := range rolesList {
					if roleMap, ok := roleItem.(map[string]any); ok {
						role := &authdto.RoleDto{
							RoleID:          getStringValue(roleMap, "id"),
							RoleName:        getStringValue(roleMap, "name"),
							RoleDescription: getStringValue(roleMap, "description"),
							ClientRole:      getBoolValue(roleMap, "clientRole"),
						}
						roles = append(roles, role)
					}
				}
			}
		}
	}

	return roles, nil
}

// getStringValue 从 map 中获取字符串值
func getStringValue(m map[string]any, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getBoolValue 从 map 中获取布尔值
func getBoolValue(m map[string]any, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// SearchUsersRequest 用户搜索请求参数
type SearchUsersRequest struct {
	First            int
	Max              int
	Username         string // 精确查询
	UsernameWildcard string // 模糊查询
	FirstName        string // 模糊查询
	Email            string // 模糊查询
	Enabled          *bool
	Exact            bool // 是否精确查询
}

// SearchUsers 搜索用户列表
func (s *RoleService) SearchUsers(params SearchUsersRequest) ([]*KeycloakUserInfoDto, error) {
	if s.svcCtx.Keycloak == nil {
		return []*KeycloakUserInfoDto{}, nil
	}

	// 构建查询参数
	searchParams := make(map[string]string)

	if params.First > 0 {
		searchParams["first"] = strconv.Itoa(params.First)
	}
	if params.Max > 0 {
		searchParams["max"] = strconv.Itoa(params.Max)
	}

	// 精确查询模式
	if params.Exact && params.Username != "" {
		searchParams["username"] = params.Username
		searchParams["exact"] = "true"
	} else {
		// 模糊查询
		if params.UsernameWildcard != "" {
			searchParams["username"] = "*" + params.UsernameWildcard + "*"
		}
		if params.FirstName != "" {
			searchParams["firstName"] = "*" + params.FirstName + "*"
		}
		if params.Email != "" {
			searchParams["email"] = "*" + params.Email + "*"
		}
	}

	if params.Enabled != nil {
		searchParams["enabled"] = strconv.FormatBool(*params.Enabled)
	}

	// 调用 Keycloak API
	users, err := s.svcCtx.Keycloak.SearchUsers(searchParams)
	if err != nil {
		s.Errorf("Failed to search users: %v", err)
		return nil, err
	}

	// 转换为 KeycloakUserInfoDto
	result := make([]*KeycloakUserInfoDto, 0, len(users))
	for i := range users {
		result = append(result, &KeycloakUserInfoDto{
			ID:                users[i].ID,
			Username:          users[i].Username,
			FirstName:         users[i].FirstName,
			LastName:          users[i].LastName,
			Email:             users[i].Email,
			Enabled:           users[i].Enabled,
			Attributes:        users[i].Attributes,
			PreferredUsername: users[i].PreferredUsername,
		})
	}

	return result, nil
}

// FetchUser 通过用户名获取单个用户
func (s *RoleService) FetchUser(username string) (*KeycloakUserInfoDto, error) {
	if s.svcCtx.Keycloak == nil {
		return nil, nil
	}

	user, err := s.svcCtx.Keycloak.FetchUser(username)
	if err != nil {
		s.Errorf("Failed to fetch user %s: %v", username, err)
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &KeycloakUserInfoDto{
		ID:                user.ID,
		Username:          user.Username,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		Enabled:           user.Enabled,
		Attributes:        user.Attributes,
		PreferredUsername: user.PreferredUsername,
	}, nil
}

// IsServiceAccount 检查用户是否为服务账户
func (s *RoleService) IsServiceAccount(attributes map[string]any) bool {
	if attributes == nil {
		return false
	}

	if v, ok := attributes["serviceAccountClientLink"]; ok {
		if str, ok := v.(string); ok && str != "" {
			return true
		}
		if arr, ok := v.([]any); ok && len(arr) > 0 {
			if str, ok := arr[0].(string); ok && str != "" {
				return true
			}
		}
	}
	return false
}

// GetStringAttribute 从属性 map 中获取字符串值
func (s *RoleService) GetStringAttribute(attributes map[string]any, key string) string {
	if attributes == nil {
		return ""
	}

	if v, ok := attributes[key]; ok {
		if str, ok := v.(string); ok && str != "" {
			return str
		}
		if arr, ok := v.([]any); ok && len(arr) > 0 {
			if str, ok := arr[0].(string); ok {
				return str
			}
		}
	}
	return ""
}

// GetIntAttribute 从属性 map 中获取整数值
func (s *RoleService) GetIntAttribute(attributes map[string]any, key string, defaultValue int) int {
	if attributes == nil {
		return defaultValue
	}

	if v, ok := attributes[key]; ok {
		if str, ok := v.(string); ok && str != "" {
			if i, err := strconv.Atoi(str); err == nil {
				return i
			}
		}
		if arr, ok := v.([]any); ok && len(arr) > 0 {
			if str, ok := arr[0].(string); ok {
				if i, err := strconv.Atoi(str); err == nil {
					return i
				}
			}
		}
	}
	return defaultValue
}

// KeycloakUserInfoDto Keycloak 用户信息
type KeycloakUserInfoDto struct {
	ID                string
	Username          string
	FirstName         string
	LastName          string
	Email             string
	Enabled           bool
	Attributes        map[string]any
	PreferredUsername string
}
