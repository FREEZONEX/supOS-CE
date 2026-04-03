package userManage

import (
	"context"

	"backend/internal/common/enums"
	"backend/internal/svc"
	"backend/internal/types"
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
	db, err := l.db()
	if err != nil {
		return nil, err
	}

	roles, err := l.listRoles(db)
	if err != nil {
		return nil, err
	}

	var (
		adminRoles []types.RoleDetail
		otherRoles []types.RoleDetail
	)
	for _, role := range roles {
		resources, denyResources, resourceErr := l.getRolePermissionLists(db, role.ID)
		if resourceErr != nil {
			return nil, resourceErr
		}
		display, _ := normalizeRoleDisplay(l.ctx, role.RoleKey, role.RoleName, role.Description)
		detail := types.RoleDetail{
			RoleID:           role.ID,
			RoleName:         firstNonEmpty(display, role.RoleName, role.RoleKey),
			ResourceList:     resources,
			DenyResourceList: denyResources,
		}
		if role.ID == enums.RoleSuperAdmin.ID || role.RoleKey == "admin" {
			adminRoles = append(adminRoles, detail)
		} else {
			otherRoles = append(otherRoles, detail)
		}
	}
	return append(adminRoles, otherRoles...), nil
}
