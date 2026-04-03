package auth

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/common/constants"
	"backend/internal/common/vo"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
)

type LoginLogic struct {
	baseAuthLogic
}

type LoginResult struct {
	Cookie *http.Cookie
	User   *vo.UserInfoVo
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		baseAuthLogic: newBaseAuthLogic(ctx, svcCtx),
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*LoginResult, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("invalid login request")
	}
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" {
		return nil, errors.Parameter.WithMsg("user.username.null")
	}
	if password == "" {
		return nil, errors.Parameter.WithMsg("user.password.null")
	}

	repo, err := l.iamRepo()
	if err != nil || repo == nil {
		return nil, errors.System.WithMsg("iam repository unavailable")
	}

	user, userErr := repo.GetUserByUsername(username)
	if userErr != nil {
		l.Errorf("load iam user failed for %s: %v", username, userErr)
		return nil, errors.System.WithMsg("failed to load user")
	}
	if user == nil || !user.Enabled {
		return nil, errors.Parameter.WithMsg("user.login.password.error")
	}

	if strings.TrimSpace(user.Password) == "" || !verifyPassword(user.Password, password) {
		return nil, errors.Parameter.WithMsg("user.login.password.error")
	}

	sessionID, err := l.createIAMSession(user.ID)
	if err != nil {
		return nil, errors.System.WithMsg("failed to create session")
	}

	currentInfo, infoErr := l.buildIAMUserInfo(user.ID)
	if infoErr != nil {
		return nil, infoErr
	}

	return &LoginResult{
		Cookie: buildSessionCookie(sessionID),
		User:   currentInfo,
	}, nil
}

func buildSessionCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     constants.AccessTokenKey,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   constants.CookieMaxAge,
		HttpOnly: false,
		Secure:   false,
	}
}
