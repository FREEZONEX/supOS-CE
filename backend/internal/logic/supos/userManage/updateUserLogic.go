package userManage

import (
	"context"
	"strings"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
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

	user, err := l.getUserEntity(strings.TrimSpace(req.UserID))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Parameter.WithMsg("user.not.exist")
	}

	db, err := l.keycloakDB()
	if err != nil {
		return nil, err
	}

	kc, err := l.keycloakClient()
	if err != nil {
		return nil, err
	}

	email := strings.TrimSpace(user.Email)
	if v := strings.TrimSpace(req.Email); v != "" && !strings.EqualFold(v, user.Email) {
		existing, err := kc.FetchUserByEmail(v)
		if err != nil {
			l.Errorf("failed to query user by email %s: %v", v, err)
			return nil, errors.System.WithMsg("failed to check email")
		}
		if existing != nil && existing.ID != user.ID {
			return nil, errors.Parameter.WithMsg("user.email.already.exists")
		}
		email = v
	}

	firstName := strings.TrimSpace(user.FirstName)
	if v := strings.TrimSpace(req.FirstName); v != "" {
		firstName = v
	}

	enabled := user.Enabled
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	attrMap, err := l.loadAttributesForUsers(db, []string{user.ID})
	if err != nil {
		return nil, err
	}
	attributes := attrMap[user.ID]
	if attributes == nil {
		attributes = make(map[string]string)
	}

	if v := strings.TrimSpace(req.Phone); v != "" {
		attributes["phone"] = v
	}
	// if v := strings.TrimSpace(req.HomePage); v != "" {
	// 	attributes["homePage"] = v
	// }
	if v := strings.TrimSpace(req.Source); v != "" {
		attributes["source"] = v
	}
	// if req.TipsEnable != nil {
	// 	attributes["tipsEnable"] = strconv.Itoa(*req.TipsEnable)
	// }
	// if req.FirstTimeLogin != nil {
	// 	attributes["firstTimeLogin"] = strconv.Itoa(*req.FirstTimeLogin)
	// }
	// if _, ok := attributes["homePage"]; !ok {
	// 	attributes["homePage"] = "/home"
	// }
	// if _, ok := attributes["tipsEnable"]; !ok {
	// 	attributes["tipsEnable"] = "1"
	// }
	// if _, ok := attributes["firstTimeLogin"]; !ok {
	// 	attributes["firstTimeLogin"] = "1"
	// }

	payload := map[string]any{
		"firstName": firstName,
		"enabled":   enabled,
		"email":     email,
	}
	payload["attributes"] = convertAttributesPayload(attributes)

	if err := kc.UpdateUser(user.ID, payload); err != nil {
		l.Errorf("failed to update user %s: %v", user.ID, err)
		return nil, errors.System.WithMsg("failed to update user")
	}

	if v := strings.TrimSpace(req.Password); v != "" {
		if err := kc.ResetPassword(user.ID, v); err != nil {
			l.Errorf("failed to reset password for user %s: %v", user.ID, err)
			return nil, errors.System.WithMsg("failed to update user password")
		}
	}

	currentRoleMap, err := l.loadRolesForUsers(db, []string{user.ID})
	if err != nil {
		return nil, err
	}
	currentRoles := currentRoleMap[user.ID]

	if len(req.RoleList) > 0 {
		if len(currentRoles) > 0 {
			if err := l.applyUserRoles(user.ID, currentRoles, false); err != nil {
				return nil, err
			}
		}
		if err := l.applyUserRoles(user.ID, req.RoleList, true); err != nil {
			return nil, err
		}
	} else if req.OperateRole != nil && *req.OperateRole {
		if len(currentRoles) > 0 {
			if err := l.applyUserRoles(user.ID, currentRoles, false); err != nil {
				return nil, err
			}
		}
	}

	l.invalidateUserCache(user.ID)
	return l.success(), nil
}
