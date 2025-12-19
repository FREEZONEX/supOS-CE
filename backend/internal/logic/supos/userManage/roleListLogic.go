package userManage

import (
	"context"
	"fmt"
	"strings"

	authdto "backend/internal/common/dto/auth"
	"backend/internal/common/enums"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/clients"

	"gitee.com/unitedrhino/share/errors"
)

type RoleListLogic struct {
	baseUserManageLogic
}

// List available roles
func NewRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleListLogic {
	return &RoleListLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *RoleListLogic) RoleList() ([]types.RoleDetail, error) {
	kc, err := l.keycloakClient()
	if err != nil {
		return nil, err
	}

	roles, err := kc.GetAllRoles()
	if err != nil {
		l.Errorf("failed to load roles from keycloak: %v", err)
		return nil, errors.System.WithMsg(fmt.Sprintf("failed to load roles: %v", err))
	}

	repo, err := l.authRepo()
	if err != nil {
		l.Errorf("init keycloak repo failed: %v", err)
		return nil, errors.System.WithMsg("failed to access keycloak repository")
	}

	roleByName := make(map[string]*clients.KeycloakRoleInfoDto, len(roles))
	allowRoles := make(map[string]*clients.KeycloakRoleInfoDto)
	orderedNames := make([]string, 0, len(roles))
	for i := range roles {
		role := &roles[i]
		roleByName[role.Name] = role
		if enums.IsIgnoredRoleID(role.ID) || enums.IsIgnoredRoleName(role.Name) || strings.HasPrefix(role.Name, "deny-") {
			continue
		}
		if _, exists := allowRoles[role.Name]; !exists {
			orderedNames = append(orderedNames, role.Name)
		}
		allowRoles[role.Name] = role
	}

	var (
		superRoles []types.RoleDetail
		otherRoles []types.RoleDetail
	)

	for _, roleName := range orderedNames {
		role := allowRoles[roleName]
		if role == nil {
			continue
		}
		var allowResources []*authdto.ResourceDto
		var denyResources []*authdto.ResourceDto
		if repo != nil {
			allowResources, err = repo.GetRoleAllowResources(l.ctx, role.ID)
			if err != nil {
				l.Errorf("load allow resources for role %s failed: %v", role.ID, err)
				return nil, errors.System.WithMsg(fmt.Sprintf("failed to load role resources: %v", err))
			}
		}
		if denyRole := roleByName[fmt.Sprintf("deny-%s", roleName)]; denyRole != nil && repo != nil {
			denyResources, err = repo.GetRoleDenyResources(l.ctx, denyRole.ID)
			if err != nil {
				l.Errorf("load deny resources for role %s failed: %v", denyRole.ID, err)
				return nil, errors.System.WithMsg(fmt.Sprintf("failed to load role resources: %v", err))
			}
		}

		displayName, desc := normalizeRoleDisplay(l.ctx,role.ID, role.Name, role.Description)
		detail := types.RoleDetail{
			RoleID:           role.ID,
			RoleName:         strings.TrimSpace(displayName),
			ResourceList:     toRoleResourceList(allowResources),
			DenyResourceList: toRoleResourceList(denyResources),
		}
		if detail.RoleName == "" {
			detail.RoleName = strings.TrimSpace(desc)
		}

		if role.ID == enums.RoleSuperAdmin.ID {
			superRoles = append(superRoles, detail)
		} else {
			otherRoles = append(otherRoles, detail)
		}
	}
	resp := append(superRoles, otherRoles...)
	return resp, nil
}
