package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	cache "backend/internal/common/cache"
	"backend/internal/common/constants"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/google/uuid"
)

type TokenLogic struct {
	baseAuthLogic
}

type TokenResult struct {
	Cookie      *http.Cookie
	RedirectURL string
}

func NewTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenLogic {
	return &TokenLogic{
		baseAuthLogic: newBaseAuthLogic(ctx, svcCtx),
	}
}

func (l *TokenLogic) Token(req *types.TokenCallbackReq) (*TokenResult, error) {
	if req == nil || strings.TrimSpace(req.Code) == "" {
		return nil, errors.Parameter.WithMsg("authorization code is required")
	}
	kc := l.keycloakClient()
	if kc == nil {
		return nil, errors.System.WithMsg("keycloak client not configured")
	}

	tokenResp, err := kc.GetKeyCloakTokenByCode(req.Code)
	if err != nil {
		l.Errorf("exchange code for token failed: %v", err)
		return nil, errors.System.WithMsg("failed to exchange token")
	}
	if tokenResp == nil || strings.TrimSpace(tokenResp.AccessToken) == "" {
		return nil, errors.System.WithMsg("empty token response from keycloak")
	}

	sessionKey := tokenResp.SessionState
	if sessionKey == "" {
		sessionKey = uuid.NewString()
	}

	raw := map[string]any{
		"access_token":       tokenResp.AccessToken,
		"refresh_token":      tokenResp.RefreshToken,
		"refresh_expires_in": tokenResp.RefreshExpiresIn,
		"expires_in":         tokenResp.ExpiresIn,
		"token_type":         tokenResp.TokenType,
		"scope":              tokenResp.Scope,
		"session_state":      sessionKey,
	}

	if cache.TokenCache != nil {
		cache.TokenCache.Set(sessionKey, &cache.TokenCacheEntry{
			Token:    tokenResp,
			Raw:      raw,
			CachedAt: time.Now(),
		})
	}

	userInfo, _, err := l.fetchUserInfo(tokenResp.AccessToken, false)
	_ = userInfo
	if err != nil {
		l.Errorf("load user info failed: %v", err)
	}

	redirect := l.svcCtx.Config.OAuthKeyCloak.SuposHome
	if userInfo != nil && strings.TrimSpace(userInfo.HomePage) != "" {
		redirect = userInfo.HomePage
	}

	cookie := &http.Cookie{
		Name:     constants.AccessTokenKey,
		Value:    sessionKey,
		Path:     "/",
		MaxAge:   constants.CookieMaxAge,
		HttpOnly: false,
		Secure:   false,
	}

	return &TokenResult{
		Cookie:      cookie,
		RedirectURL: redirect,
	}, nil
}
