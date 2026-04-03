package auth

import (
	"context"
	"strings"

	"backend/internal/svc"
)

type LogoutLogic struct {
	baseAuthLogic
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		baseAuthLogic: newBaseAuthLogic(ctx, svcCtx),
	}
}

func (l *LogoutLogic) Logout(sessionKey string) error {
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		return nil
	}

	if session, err := l.getIAMSession(sessionKey); err == nil && session != nil {
		if revokeErr := l.removeIAMSession(session.ID); revokeErr != nil {
			l.Errorf("revoke iam session failed: %v", revokeErr)
		}
		return nil
	}
	l.removeSession(sessionKey)
	return nil
}
