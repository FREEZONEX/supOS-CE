package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type ResetPasswordLogic struct {
	baseUserManageLogic
}

// Reset user password by admin
func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.AdminResetPwdReq) (*types.OperationResult, error) {
	if req == nil || strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, errors.Parameter.WithMsg("userId and password are required")
	}

	kc, err := l.keycloakClient()
	if err != nil {
		return nil, err
	}

	if err := kc.ResetPassword(strings.TrimSpace(req.UserID), strings.TrimSpace(req.Password)); err != nil {
		l.Errorf("failed to reset password for user %s: %v", req.UserID, err)
		return nil, errors.System.WithMsg("failed to reset password")
	}

	l.invalidateUserCache(strings.TrimSpace(req.UserID))
	return l.success(), nil
}
