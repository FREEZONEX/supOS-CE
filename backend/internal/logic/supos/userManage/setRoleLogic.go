package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gorm.io/gorm"
)

type SetRoleLogic struct {
	baseUserManageLogic
}

// Assign roles to user
func NewSetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetRoleLogic {
	return &SetRoleLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *SetRoleLogic) SetRole(req *types.UserSetRoleReq) (*types.OperationResult, error) {
	if req == nil || strings.TrimSpace(req.UserID) == "" {
		return nil, errors.Parameter.WithMsg("userId is required")
	}

	db, err := l.db()
	if err != nil {
		return nil, err
	}
	user, err := l.getIAMUser(db, req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Parameter.WithMsg("user.not.exist")
	}

	roleIDs, err := l.resolveRoleIDs(db, req.RoleList)
	if err != nil {
		return nil, err
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if req.Type == 2 {
			return l.removeUserRoles(tx, user.ID, roleIDs)
		}
		return l.addUserRoles(tx, user.ID, roleIDs)
	}); err != nil {
		return nil, err
	}

	l.invalidateUserCache(user.ID)
	return l.success(), nil
}
