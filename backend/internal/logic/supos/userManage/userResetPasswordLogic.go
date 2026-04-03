package userManage

import (
	"context"
	"strings"

	authlogic "backend/internal/logic/supos/auth"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gorm.io/gorm"
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

	db, err := l.db()
	if err != nil {
		return nil, err
	}

	user, err := l.getIAMUser(db, req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || !strings.EqualFold(user.Username, req.Username) {
		return nil, errors.Parameter.WithMsg("user.not.exist")
	}

	if strings.TrimSpace(user.Password) == "" || !authlogic.VerifyPassword(user.Password, req.Password) {
		return nil, errors.Parameter.WithMsg("user.login.password.error")
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return l.upsertUserPassword(tx, user.ID, req.NewPassword)
	}); err != nil {
		return nil, err
	}

	l.invalidateUserCache(user.ID)
	return l.success(), nil
}
