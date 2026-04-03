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
)

type UpdateUserLogic struct {
	baseUserManageLogic
}

// Update user profile
func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UserUpdateReq) (*types.OperationResult, error) {
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

	if username := strings.TrimSpace(req.Username); username != "" && !strings.EqualFold(username, user.Username) {
		if existing, err := l.getIAMUserByUsername(db, username); err != nil {
			return nil, err
		} else if existing != nil && existing.ID != user.ID {
			return nil, errors.Parameter.WithMsg("user.username.already.exists")
		}
		user.Username = username
	}

	if email := strings.TrimSpace(req.Email); email != "" && !strings.EqualFold(email, user.Email) {
		if existing, err := l.getIAMUserByEmail(db, email); err != nil {
			return nil, err
		} else if existing != nil && existing.ID != user.ID {
			return nil, errors.Parameter.WithMsg("user.email.already.exists")
		}
		user.Email = email
	}

	if firstName := strings.TrimSpace(req.FirstName); firstName != "" {
		user.DisplayName = firstName
	}
	if req.Enabled != nil {
		user.Enabled = *req.Enabled
	}
	if source := strings.TrimSpace(req.Source); source != "" {
		user.Source = source
	}
	now := time.Now()
	if phone := strings.TrimSpace(req.Phone); phone != "" {
		user.Phone = phone
	}
	if homePage := strings.TrimSpace(req.HomePage); homePage != "" {
		user.HomePage = homePage
	}
	if req.FirstTimeLogin != nil {
		user.FirstTimeLogin = *req.FirstTimeLogin
	}
	if req.TipsEnable != nil {
		user.TipsEnable = *req.TipsEnable
	}

	roleIDs := []string(nil)
	if len(req.RoleList) > 0 {
		roleIDs, err = l.resolveRoleIDs(db, req.RoleList)
		if err != nil {
			return nil, err
		}
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if txErr := tx.Model(&relationDB.IamUser{}).Where("id = ?", user.ID).Updates(map[string]any{
			"username":         user.Username,
			"display_name":     firstNonEmpty(user.DisplayName, user.Username),
			"email":            user.Email,
			"enabled":          user.Enabled,
			"source":           user.Source,
			"phone":            user.Phone,
			"home_page":        user.HomePage,
			"first_time_login": user.FirstTimeLogin,
			"tips_enable":      user.TipsEnable,
			"updated_at":       now,
		}).Error; txErr != nil {
			l.Errorf("update user failed: %v", txErr)
			return errors.System.WithMsg("failed to update user")
		}
		if password := strings.TrimSpace(req.Password); password != "" {
			if passwordErr := l.upsertUserPassword(tx, user.ID, password); passwordErr != nil {
				return passwordErr
			}
		}
		if len(req.RoleList) > 0 {
			if roleErr := l.replaceUserRoles(tx, user.ID, roleIDs); roleErr != nil {
				return roleErr
			}
		} else if req.OperateRole != nil && *req.OperateRole {
			if roleErr := l.replaceUserRoles(tx, user.ID, nil); roleErr != nil {
				return roleErr
			}
		}
		if !user.Enabled {
			if revokeErr := l.revokeUserSessions(tx, user.ID); revokeErr != nil {
				return revokeErr
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	l.invalidateUserCache(user.ID)
	return l.success(), nil
}
