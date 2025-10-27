package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
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

	add := true
	if req.Type == 2 {
		add = false
	}

	if err := l.applyUserRoles(strings.TrimSpace(req.UserID), req.RoleList, add); err != nil {
		return nil, err
	}

	l.invalidateUserCache(strings.TrimSpace(req.UserID))
	return l.success(), nil
}
