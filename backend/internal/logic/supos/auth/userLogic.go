package auth

import (
	"context"
	"strings"

	"backend/internal/common/vo"
	"backend/internal/svc"
)

type UserLogic struct {
	baseAuthLogic
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		baseAuthLogic: newBaseAuthLogic(ctx, svcCtx),
	}
}

func (l *UserLogic) User(sessionKey string) (any, error) {
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		if l.authDisabled() {
			return vo.Guest(), nil
		}
		return "not found user info", nil
	}

	if session, err := l.getIAMSession(sessionKey); err == nil && session != nil {
		if !l.isIAMSessionActive(session) {
			l.removeSession(sessionKey)
			if l.authDisabled() {
				return vo.Guest(), nil
			}
			return "not found user info", nil
		}
		userInfo, infoErr := l.buildIAMUserInfo(session.UserID)
		if infoErr == nil && userInfo != nil {
			_ = l.touchIAMSession(sessionKey)
			userInfo.SuperAdmin = userInfo.IsSuperAdmin()
			return userInfo, nil
		}
		l.Errorf("build iam user info failed: %v", infoErr)
	}

	if l.authDisabled() {
		return vo.Guest(), nil
	}
	l.removeSession(sessionKey)
	return "not found user info", nil
}
