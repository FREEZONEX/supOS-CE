package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type DeleteUserLogic struct {
	baseUserManageLogic
}

// Remove user by id
func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.UserDeleteReq) (*types.OperationResult, error) {
	if req == nil || strings.TrimSpace(req.ID) == "" {
		return nil, errors.Parameter.WithMsg("userId is required")
	}

	kc, err := l.keycloakClient()
	if err != nil {
		return nil, err
	}

	if err := kc.DeleteUser(strings.TrimSpace(req.ID)); err != nil {
		l.Errorf("failed to delete user %s: %v", req.ID, err)
		return nil, errors.System.WithMsg("failed to delete user")
	}

	l.invalidateUserCache(strings.TrimSpace(req.ID))
	return l.success(), nil
}
