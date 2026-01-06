package userManage

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/common/constants"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type CreateUserLogic struct {
	baseUserManageLogic
}

// Create user
func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *CreateUserLogic) CreateUser(req *types.UserCreateReq) (*types.OperationResult, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("request body is empty")
	}
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, errors.Parameter.WithMsg("user.username.empty")
	}
	if len(username) < 3 || len(username) > 30 {
		return nil, errors.Parameter.WithMsg("user.username.invalid")
	}
	password := strings.TrimSpace(req.Password)
	if password == "" {
		return nil, errors.Parameter.WithMsg("user.password.empty")
	}

	kc, err := l.keycloakClient()
	if err != nil {
		return nil, err
	}

	existing, err := kc.FetchUser(username)
	if err != nil {
		l.Errorf("failed to query user by username %s: %v", username, err)
		return nil, errors.System.WithMsg("failed to check username")
	}
	if existing != nil {
		return nil, errors.Parameter.WithMsg("user.username.already.exists")
	}

	email := strings.TrimSpace(req.Email)
	if email != "" {
		existing, err := kc.FetchUserByEmail(email)
		if err != nil {
			l.Errorf("failed to query user by email %s: %v", email, err)
			return nil, errors.System.WithMsg("failed to check email")
		}
		if existing != nil {
			return nil, errors.Parameter.WithMsg("user.email.already.exists")
		}
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	attributes := map[string]any{
		"firstTimeLogin": []string{"1"},
		"tipsEnable":     []string{"1"},
		"homePage":       []string{constants.DefaultHomepage},
	}
	if phone := strings.TrimSpace(req.Phone); phone != "" {
		attributes["phone"] = []string{phone}
	}
	if source := strings.TrimSpace(req.Source); source != "" {
		attributes["source"] = []string{source}
	}

	userPayload := map[string]any{
		"username":   username,
		"enabled":    enabled,
		"email":      email,
		"firstName":  strings.TrimSpace(req.FirstName),
		"attributes": attributes,
	}

	var createdUserID string
	var cleanup bool
	defer func() {
		if !cleanup || createdUserID == "" {
			return
		}
		if err := kc.DeleteUser(createdUserID); err != nil {
			l.Errorf("failed to cleanup user %s: %v", createdUserID, err)
		}
	}()

	userID, err := kc.CreateUser(userPayload)
	if err != nil {
		l.Errorf("failed to create keycloak user %s: %v", username, err)
		return nil, errors.System.WithMsg(fmt.Sprintf("failed to create user: %v", err))
	}
	createdUserID = userID
	cleanup = true

	if err := kc.ResetPassword(userID, password); err != nil {
		l.Errorf("failed to reset password for user %s: %v", username, err)
		return nil, errors.System.WithMsg("failed to set user password")
	}

	if len(req.RoleList) > 0 {
		if err := l.applyUserRoles(userID, req.RoleList, true); err != nil {
			l.Errorf("failed to assign roles for user %s: %v", username, err)
			return nil, err
		}
	}

	cleanup = false
	l.invalidateUserCache(userID)
	return l.success(), nil
}
