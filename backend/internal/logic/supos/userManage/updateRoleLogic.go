package userManage

import (
	"context"
	"strings"
	"time"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gorm.io/gorm"
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

	db, err := l.db()
	if err != nil {
		return nil, err
	}

	role, err := l.getRoleByID(db, req.ID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.Parameter.WithMsg("role.no.exist")
	}
	if l.protectedRole(role) {
		return nil, errors.Parameter.WithMsg("role.super.update")
	}

	newName := strings.TrimSpace(req.Name)
	if !strings.EqualFold(role.RoleName, newName) {
		if existing, err := l.getRoleByName(db, newName); err != nil {
			return nil, err
		} else if existing != nil && existing.ID != role.ID {
			return nil, errors.Parameter.WithMsg("role.name.exist")
		}
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"role_name":   newName,
			"role_key":    normalizeRoleKey(newName),
			"description": newName,
			"updated_at":  time.Now(),
		}
		if txErr := tx.Model(&relationDB.IamRole{}).Where("id = ?", role.ID).Updates(updates).Error; txErr != nil {
			l.Errorf("update role failed: %v", txErr)
			return errors.System.WithMsg("failed to update role")
		}
		return l.replaceRoleResources(tx, role.ID, req.AllowResourceList, req.DenyResourceList)
	}); err != nil {
		return nil, err
	}

	resourceList, denyResourceList, err := l.getRolePermissionLists(db, role.ID)
	if err != nil {
		return nil, err
	}
	display, _ := normalizeRoleDisplay(l.ctx, normalizeRoleKey(newName), newName, newName)
	return &types.RoleDetail{
		RoleID:           role.ID,
		RoleName:         firstNonEmpty(display, newName, role.RoleKey),
		ResourceList:     resourceList,
		DenyResourceList: denyResourceList,
	}, nil
}
