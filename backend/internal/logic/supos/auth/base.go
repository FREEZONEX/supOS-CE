package auth

import (
	"context"
	"os"
	"strings"

	authsvc "backend/internal/common/authsvc"
	cache "backend/internal/common/cache"
	"backend/internal/common/constants"
	"backend/internal/common/vo"
	"backend/internal/svc"
	"backend/share/clients"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type baseAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func newBaseAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) baseAuthLogic {
	return baseAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *baseAuthLogic) getTokenEntry(token string) (*cache.TokenCacheEntry, bool) {
	if token == "" || cache.TokenCache == nil {
		return nil, false
	}
	return cache.TokenCache.Get(token)
}

func (l *baseAuthLogic) removeSession(sessionKey string) {
	if sessionKey == "" {
		return
	}
	if kc := l.keycloakClient(); kc != nil {
		if err := kc.RemoveSession(sessionKey); err != nil {
			l.Errorf("remove session failed: %v", err)
		}
	}
	if cache.TokenCache != nil {
		cache.TokenCache.Delete(sessionKey)
	}
}

func (l *baseAuthLogic) keycloakClient() *clients.KeycloakClient {
	return l.svcCtx.Keycloak
}

func (l *baseAuthLogic) authDisabled() bool {
	return strings.EqualFold(os.Getenv("SYS_OS_AUTH_ENABLE"), "false") || os.Getenv("SYS_OS_AUTH_ENABLE") == ""
}

func (l *baseAuthLogic) fetchUserInfo(accessToken string, allowCache bool) (*vo.UserInfoVo, string, error) {
	kc := l.keycloakClient()
	if kc == nil {
		return nil, "", errors.System.WithMsg("keycloak client not configured")
	}
	return authsvc.FetchUserInfo(l.ctx, kc, accessToken, allowCache, constants.DefaultHomepage, l.svcCtx.Config.OAuthKeyCloak.Realm)
}
