package userManage

import (
	"context"
	"strings"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gorm.io/gorm"
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

	db, err := l.db()
	if err != nil {
		return nil, err
	}

	user, err := l.getIAMUser(db, req.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return l.success(), nil
	}
	if strings.EqualFold(user.Username, "tier0") {
		return nil, errors.Parameter.WithMsg("user.supos.delete")
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if txErr := tx.Where("user_id = ?", user.ID).Delete(&relationDB.IamSession{}).Error; txErr != nil {
			l.Errorf("delete user sessions failed: %v", txErr)
			return errors.System.WithMsg("failed to delete user")
		}
		if txErr := tx.Where("user_id = ?", user.ID).Delete(&relationDB.IamUserRole{}).Error; txErr != nil {
			l.Errorf("delete user roles failed: %v", txErr)
			return errors.System.WithMsg("failed to delete user")
		}
		if txErr := tx.Where("id = ?", user.ID).Delete(&relationDB.IamUser{}).Error; txErr != nil {
			l.Errorf("delete user failed: %v", txErr)
			return errors.System.WithMsg("failed to delete user")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	l.invalidateUserCache(user.ID)
	return l.success(), nil
}
