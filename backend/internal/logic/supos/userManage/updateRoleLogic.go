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

type UpdateRoleLogic struct {
	baseUserManageLogic
}

// Update an existing role
func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.RoleSaveReq) (*types.RoleDetail, error) {
	if req == nil || strings.TrimSpace(req.ID) == "" {
		return nil, errors.Parameter.WithMsg("role.id.empty")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.Parameter.WithMsg("role.name.empty")
	}

	kc, err := l.keycloakClient()
	if err != nil {
		return nil, err
	}

	role, err := kc.FetchRoleByID(strings.TrimSpace(req.ID))
	if err != nil {
		l.Errorf("failed to fetch role by id %s: %v", req.ID, err)
		return nil, errors.System.WithMsg("failed to load role")
	}
	if role == nil {
		return nil, errors.Parameter.WithMsg("role.no.exist")
	}

	if parsed, ok := enums.RoleParse(role.ID); ok {
		if parsed.ID == enums.RoleSuperAdmin.ID || parsed.ID == enums.RoleNormalUser.ID {
			return nil, errors.Parameter.WithMsg("role.super.update")
		}
	}

	roleName := role.Name
	denyRoleName := fmt.Sprintf("deny-%s", roleName)
	denyRole, err := kc.FetchRole(denyRoleName)
	if err != nil {
		l.Errorf("failed to fetch deny role %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to load deny role")
	}

	allowURIs := dedupeStrings(append(stringSlice(req.AllowResourceList), enums.DefaultAllowURIs...))
	denyURIs := dedupeStrings(stringSlice(req.DenyResourceList))

	if err := l.updateResource(fmt.Sprintf("%s-resource", roleName), allowURIs); err != nil {
		return nil, err
	}
	if denyRole != nil {
		if err := l.updateResource(fmt.Sprintf("%s-resource", denyRoleName), denyURIs); err != nil {
			return nil, err
		}
	}

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
	if repo != nil && denyRole != nil {
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

	return detail, nil
}

func (l *UpdateRoleLogic) updateResource(resourceName string, uris []string) error {
	kc, err := l.keycloakClient()
	if err != nil {
		return err
	}
	resource, err := kc.FetchResource(resourceName)
	if err != nil {
		l.Errorf("failed to fetch resource %s: %v", resourceName, err)
		return errors.System.WithMsg("failed to load role resource")
	}
	if resource == nil {
		l.Errorf("resource %s not found when updating role", resourceName)
		return errors.Parameter.WithMsg("role.resource.not.exist")
	}

	payload := map[string]any{
		"name":        resource.Name,
		"displayName": resource.DisplayName,
		"type":        resource.Type,
		"uris":        uris,
	}

	if err := kc.UpdateResource(resource.ID, payload); err != nil {
		l.Errorf("failed to update resource %s: %v", resourceName, err)
		return errors.System.WithMsg("failed to update role resource")
	}
	return nil
}
