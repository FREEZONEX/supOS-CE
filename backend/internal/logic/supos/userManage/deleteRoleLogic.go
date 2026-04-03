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
	"gorm.io/gorm/clause"
)

type DeleteRoleLogic struct {
	baseUserManageLogic
}

// Delete role by id
func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.RoleIDPathReq) (*types.OperationResult, error) {
	if req == nil || strings.TrimSpace(req.ID) == "" {
		return nil, errors.Parameter.WithMsg("role.id.empty")
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
		return l.success(), nil
	}
	if l.protectedRole(role) {
		return nil, errors.Parameter.WithMsg("role.super.delete")
	}

	affectedUsers := make([]string, 0)
	if err := db.Transaction(func(tx *gorm.DB) error {
		if txErr := tx.Model(&relationDB.IamUserRole{}).
			Where("role_id = ?", role.ID).
			Distinct("user_id").
			Pluck("user_id", &affectedUsers).Error; txErr != nil {
			l.Errorf("load affected users failed: %v", txErr)
			return errors.System.WithMsg("failed to delete role")
		}
		if txErr := tx.Where("role_id = ?", role.ID).Delete(&relationDB.IamUserRole{}).Error; txErr != nil {
			l.Errorf("delete user-role mappings failed: %v", txErr)
			return errors.System.WithMsg("failed to delete role")
		}
		if txErr := tx.Where("role_id = ?", role.ID).Delete(&relationDB.IamRoleResource{}).Error; txErr != nil {
			l.Errorf("delete role resources failed: %v", txErr)
			return errors.System.WithMsg("failed to delete role")
		}
		if txErr := tx.Where("id = ?", role.ID).Delete(&relationDB.IamRole{}).Error; txErr != nil {
			l.Errorf("delete role failed: %v", txErr)
			return errors.System.WithMsg("failed to delete role")
		}

		if len(affectedUsers) == 0 {
			return nil
		}
		defaultRoleID, roleErr := l.getDefaultUserRoleID(tx)
		if roleErr != nil {
			return roleErr
		}
		now := time.Now()
		operatorID := l.operatorID()
		for _, userID := range affectedUsers {
			var remaining int64
			if txErr := tx.Model(&relationDB.IamUserRole{}).Where("user_id = ?", userID).Count(&remaining).Error; txErr != nil {
				l.Errorf("count remaining roles failed: %v", txErr)
				return errors.System.WithMsg("failed to delete role")
			}
			if remaining > 0 {
				continue
			}
			item := relationDB.IamUserRole{
				UserID:    userID,
				RoleID:    defaultRoleID,
				CreatedAt: now,
				CreatedBy: operatorID,
			}
			if txErr := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&item).Error; txErr != nil {
				l.Errorf("assign default role failed: %v", txErr)
				return errors.System.WithMsg("failed to delete role")
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	for _, userID := range affectedUsers {
		l.invalidateUserCache(userID)
	}
	return l.success(), nil
}
