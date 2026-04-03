package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gorm.io/gorm"
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

	if err := db.Transaction(func(tx *gorm.DB) error {
		return l.upsertUserPassword(tx, user.ID, req.Password)
	}); err != nil {
		return nil, err
	}

	l.invalidateUserCache(user.ID)
	return l.success(), nil
}
