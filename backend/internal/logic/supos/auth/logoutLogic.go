package auth

import (
	"context"
	"strings"

	authsvc "backend/internal/common/authsvc"
	cache "backend/internal/common/cache"
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

	entry, ok := l.getTokenEntry(sessionKey)
	if ok && entry != nil && entry.Token != nil && entry.Token.RefreshToken != "" {
		if kc := l.keycloakClient(); kc != nil {
			if err := kc.Logout(entry.Token.RefreshToken); err != nil {
				l.Errorf("logout keycloak session failed: %v", err)
				return err
			}
		}
	}
	if entry != nil && entry.Token != nil {
		if claims, err := authsvc.DecodeJWTClaims(entry.Token.AccessToken); err == nil {
			if sub := authsvc.ClaimString(claims, "sub"); sub != "" && cache.UserInfoCache != nil {
				cache.UserInfoCache.Delete(sub)
			}
		}
	}
	l.removeSession(sessionKey)
	return nil
}
