package auth

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/common/vo"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/errors"
)

type UserInfoLogic struct {
	baseAuthLogic
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		baseAuthLogic: newBaseAuthLogic(ctx, svcCtx),
	}
}

func (l *UserInfoLogic) UserInfo(sessionKey string) (*vo.UserInfoVo, *http.Cookie, error) {
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		return nil, nil, errors.NotLogin.WithMsg("cookie missing")
	}

	if session, err := l.getIAMSession(sessionKey); err == nil && session != nil {
		if !l.isIAMSessionActive(session) {
			l.removeSession(sessionKey)
			return nil, nil, errors.NotLogin.WithMsg("token expired")
		}
		userInfo, infoErr := l.buildIAMUserInfo(session.UserID)
		if infoErr == nil && userInfo != nil {
			if touchErr := l.touchIAMSession(sessionKey); touchErr != nil {
				l.Errorf("touch iam session failed: %v", touchErr)
			}
			return userInfo, buildSessionCookie(sessionKey), nil
		}
		l.Errorf("build iam user info failed: %v", infoErr)
	}

	l.removeSession(sessionKey)
	return nil, nil, errors.NotLogin.WithMsg("token not found")
}
