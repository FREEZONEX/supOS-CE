package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type UserResetPasswordLogic struct {
	baseUserManageLogic
}

// Reset password by user
func NewUserResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserResetPasswordLogic {
	return &UserResetPasswordLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *UserResetPasswordLogic) UserResetPassword(req *types.UserResetPwdReq) (*types.OperationResult, error) {
	if req == nil || strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Username) == "" ||
		strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.NewPassword) == "" {
		return nil, errors.Parameter.WithMsg("userId, username and passwords are required")
	}

	kc, err := l.keycloakClient()
	if err != nil {
		return nil, err
	}

	if _, err := kc.Login(strings.TrimSpace(req.Username), strings.TrimSpace(req.Password)); err != nil {
		return nil, errors.Parameter.WithMsg("user.login.password.error")
	}

	if err := kc.ResetPassword(strings.TrimSpace(req.UserID), strings.TrimSpace(req.NewPassword)); err != nil {
		l.Errorf("failed to reset password for user %s: %v", req.UserID, err)
		return nil, errors.System.WithMsg("failed to reset password")
	}

	l.invalidateUserCache(strings.TrimSpace(req.UserID))
	return l.success(), nil
}
