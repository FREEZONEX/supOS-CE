package userManage

import (
	"context"
	"fmt"
	"strings"

	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/enums"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type CreateRoleLogic struct {
	baseUserManageLogic
}

// Create a new role
func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.RoleSaveReq) (*types.RoleDetail, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" {
		return nil, errors.Parameter.WithMsg("role.name.empty")
	}
	roleName := strings.TrimSpace(req.Name)
	denyRoleName := fmt.Sprintf("deny-%s", roleName)

	for _, r := range enums.AllRoles {
		if strings.EqualFold(roleName, r.Name) || strings.EqualFold(roleName, r.Comment) {
			return nil, errors.Parameter.WithMsg("role.name.exist")
		}
	}

	kc, err := l.keycloakClient()
	if err != nil {
		return nil, err
	}

	allRoles, err := kc.GetAllRoles()
	if err != nil {
		l.Errorf("failed to load roles: %v", err)
		return nil, errors.System.WithMsg("failed to load roles")
	}

	roleCount := 0
	for _, role := range allRoles {
		if enums.IsIgnoredRoleID(role.ID) || strings.HasPrefix(role.Name, "deny-") {
			continue
		}
		roleCount++
	}
	if roleCount >= 10 {
		return nil, errors.Parameter.WithMsg("role.max.limit")
	}

	if existing, err := kc.FetchRole(roleName); err != nil {
		l.Errorf("failed to query role by name %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to query role")
	} else if existing != nil {
		return nil, errors.Parameter.WithMsg("role.name.exist")
	}

	var cleanup []func()
	success := false
	defer func() {
		if success {
			return
		}
		for i := len(cleanup) - 1; i >= 0; i-- {
			func(f func()) {
				defer func() {
					if r := recover(); r != nil {
						l.Errorf("cleanup panic: %v", r)
					}
				}()
				f()
			}(cleanup[i])
		}
	}()

	if err := kc.CreateRole(roleName, roleName); err != nil {
		l.Errorf("failed to create role %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to create role")
	}
	cleanup = append(cleanup, func() { _ = kc.DeleteRole(roleName) })

	role, err := kc.FetchRole(roleName)
	if err != nil || role == nil {
		l.Errorf("failed to fetch created role %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to load created role")
	}

	if err := kc.CreateRole(denyRoleName, denyRoleName); err != nil {
		l.Errorf("failed to create deny role %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to create deny role")
	}
	cleanup = append(cleanup, func() { _ = kc.DeleteRole(denyRoleName) })

	denyRole, err := kc.FetchRole(denyRoleName)
	if err != nil || denyRole == nil {
		l.Errorf("failed to fetch deny role %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to load deny role")
	}

	allowURIs := dedupeStrings(append(stringSlice(req.AllowResourceList), enums.DefaultAllowURIs...))
	denyURIs := dedupeStrings(stringSlice(req.DenyResourceList))

	allowResource, err := kc.CreateResource(fmt.Sprintf("%s-resource", roleName), "URL", allowURIs)
	if err != nil {
		l.Errorf("failed to create allow resource for role %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to save role resource")
	}
	cleanup = append(cleanup, func() { _ = kc.DeleteResource(fmt.Sprintf("%s-resource", roleName)) })

	denyResource, err := kc.CreateResource(fmt.Sprintf("%s-resource", denyRoleName), "URL", denyURIs)
	if err != nil {
		l.Errorf("failed to create deny resource for role %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to save deny resource")
	}
	cleanup = append(cleanup, func() { _ = kc.DeleteResource(fmt.Sprintf("%s-resource", denyRoleName)) })

	allowPolicyName := fmt.Sprintf("%s-policy", roleName)
	allowPolicy, err := kc.CreatePolicy(allowPolicyName, allowPolicyName, role.ID)
	if err != nil {
		l.Errorf("failed to create allow policy for role %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to save role policy")
	}
	cleanup = append(cleanup, func() { _ = kc.DeletePolicy(allowPolicyName) })

	denyPolicyName := fmt.Sprintf("%s-policy", denyRoleName)
	denyPolicy, err := kc.CreatePolicy(denyPolicyName, denyPolicyName, denyRole.ID)
	if err != nil {
		l.Errorf("failed to create deny policy for role %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to save deny policy")
	}
	cleanup = append(cleanup, func() { _ = kc.DeletePolicy(denyPolicyName) })

	allowPermissionName := fmt.Sprintf("%s-permission", roleName)
	if err := kc.CreatePermission(allowPermissionName, allowPermissionName, allowPolicy.ID, allowResource.ID); err != nil {
		l.Errorf("failed to create permission for role %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to save role permission")
	}
	cleanup = append(cleanup, func() { _ = kc.DeletePermission(allowPermissionName) })

	denyPermissionName := fmt.Sprintf("%s-permission", denyRoleName)
	if err := kc.CreatePermission(denyPermissionName, denyPermissionName, denyPolicy.ID, denyResource.ID); err != nil {
		l.Errorf("failed to create deny permission for role %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to save deny permission")
	}
	cleanup = append(cleanup, func() { _ = kc.DeletePermission(denyPermissionName) })

	repo, err := l.authRepo()
	if err != nil {
		l.Errorf("init keycloak repo failed: %v", err)
		return nil, errors.System.WithMsg("failed to access keycloak repository")
	}

	var allowResourceList []*authdto.ResourceDto
	if repo != nil {
		allowResourceList, err = repo.GetRoleAllowResources(l.ctx, role.ID)
		if err != nil {
			l.Errorf("load allow resources for role %s failed: %v", role.ID, err)
			return nil, errors.System.WithMsg("failed to load role resources")
		}
	}
	if allowResourceList == nil {
		allowResourceList = dtoFromURIs(allowURIs)
	}

	var denyResourceList []*authdto.ResourceDto
	if repo != nil {
		denyResourceList, err = repo.GetRoleDenyResources(l.ctx, denyRole.ID)
		if err != nil {
			l.Errorf("load deny resources for role %s failed: %v", denyRole.ID, err)
			return nil, errors.System.WithMsg("failed to load role resources")
		}
	}
	if denyResourceList == nil {
		denyResourceList = dtoFromURIs(denyURIs)
	}

	displayName, desc := normalizeRoleDisplay(l.ctx, role.ID, role.Name, role.Description)
	detail := &types.RoleDetail{
		RoleID:           role.ID,
		RoleName:         strings.TrimSpace(displayName),
		ResourceList:     toRoleResourceList(allowResourceList),
		DenyResourceList: toRoleResourceList(denyResourceList),
	}
	if detail.RoleName == "" {
		detail.RoleName = strings.TrimSpace(desc)
	}

	success = true
	return detail, nil
}

func dtoFromURIs(uris []string) []*authdto.ResourceDto {
	if len(uris) == 0 {
		return nil
	}
	result := make([]*authdto.ResourceDto, 0, len(uris))
	for _, uri := range uris {
		uri = strings.TrimSpace(uri)
		if uri == "" {
			continue
		}
		result = append(result, &authdto.ResourceDto{URI: uri})
	}
	return result
}
