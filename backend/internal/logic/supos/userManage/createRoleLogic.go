package userManage

import (
	"context"
	"strings"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	if normalizeBuiltinRoleKey(req.Name) != "" {
		return nil, errors.Parameter.WithMsg("role.name.exist")
	}

	db, err := l.db()
	if err != nil {
		return nil, err
	}

	roleCount, err := l.countActiveRoles(db)
	if err != nil {
		return nil, err
	}
	if roleCount >= 10 {
		return nil, errors.Parameter.WithMsg("role.max.limit")
	}

	if existing, err := l.getRoleByName(db, req.Name); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, errors.Parameter.WithMsg("role.name.exist")
	}

	role := relationDB.IamRole{
		ID:          uuid.NewString(),
		RoleKey:     normalizeRoleKey(req.Name),
		RoleName:    strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Name),
		Builtin:     false,
		Status:      1,
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if txErr := tx.Create(&role).Error; txErr != nil {
			l.Errorf("create role failed: %v", txErr)
			return errors.System.WithMsg("failed to create role")
		}
		return l.replaceRoleResources(tx, role.ID, req.AllowResourceList, req.DenyResourceList)
	}); err != nil {
		return nil, err
	}

	resourceList, denyResourceList, err := l.getRolePermissionLists(db, role.ID)
	if err != nil {
		return nil, err
	}
	display, _ := normalizeRoleDisplay(l.ctx, role.RoleKey, role.RoleName, role.Description)
	return &types.RoleDetail{
		RoleID:           role.ID,
		RoleName:         firstNonEmpty(display, role.RoleName, role.RoleKey),
		ResourceList:     resourceList,
		DenyResourceList: denyResourceList,
	}, nil
}
