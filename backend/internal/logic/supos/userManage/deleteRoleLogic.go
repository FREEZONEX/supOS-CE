package userManage

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/common/enums"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type DeleteRoleLogic struct {
	baseUserManageLogic
}

// Delete role by id
func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.RoleIDPathReq) (*types.OperationResult, error) {
	if req == nil || strings.TrimSpace(req.ID) == "" {
		return nil, errors.Parameter.WithMsg("role.id.empty")
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
		return l.success(), nil
	}

	if parsed, ok := enums.RoleParse(role.ID); ok {
		if parsed.ID == enums.RoleSuperAdmin.ID || parsed.ID == enums.RoleNormalUser.ID {
			return nil, errors.Parameter.WithMsg("role.super.delete")
		}
	}

	roleName := role.Name
	denyRoleName := fmt.Sprintf("deny-%s", roleName)

	if err := kc.DeletePermission(fmt.Sprintf("%s-permission", roleName)); err != nil {
		l.Errorf("failed to delete permission %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to delete role permission")
	}
	if err := kc.DeletePermission(fmt.Sprintf("%s-permission", denyRoleName)); err != nil {
		l.Errorf("failed to delete deny permission %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to delete deny role permission")
	}

	if err := kc.DeletePolicy(fmt.Sprintf("%s-policy", roleName)); err != nil {
		l.Errorf("failed to delete policy %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to delete role policy")
	}
	if err := kc.DeletePolicy(fmt.Sprintf("%s-policy", denyRoleName)); err != nil {
		l.Errorf("failed to delete deny policy %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to delete deny role policy")
	}

	if err := kc.DeleteResource(fmt.Sprintf("%s-resource", roleName)); err != nil {
		l.Errorf("failed to delete resource %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to delete role resource")
	}
	if err := kc.DeleteResource(fmt.Sprintf("%s-resource", denyRoleName)); err != nil {
		l.Errorf("failed to delete deny resource %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to delete deny role resource")
	}

	if err := kc.DeleteRole(roleName); err != nil {
		l.Errorf("failed to delete role %s: %v", roleName, err)
		return nil, errors.System.WithMsg("failed to delete role")
	}
	if err := kc.DeleteRole(denyRoleName); err != nil {
		l.Errorf("failed to delete deny role %s: %v", denyRoleName, err)
		return nil, errors.System.WithMsg("failed to delete deny role")
	}

	return l.success(), nil
}
