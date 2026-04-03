package oauth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
)

type TokenLogic struct {
	baseOAuthLogic
}

func NewTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenLogic {
	return &TokenLogic{
		baseOAuthLogic: newBaseOAuthLogic(ctx, svcCtx),
	}
}

func (l *TokenLogic) Token(req *types.OAuthTokenReq) (*types.OAuthTokenResp, error) {
	if req == nil {
		return nil, newOAuthError(http.StatusBadRequest, "invalid_request", "missing token request")
	}
	if grantType := strings.TrimSpace(req.GrantType); grantType != "" && !strings.EqualFold(grantType, "authorization_code") {
		return nil, newOAuthError(http.StatusBadRequest, "unsupported_grant_type", "only authorization_code is supported")
	}
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.ClientID) == "" {
		return nil, newOAuthError(http.StatusBadRequest, "invalid_request", "code and client_id are required")
	}

	repo, err := l.repo()
	if err != nil || repo == nil {
		l.Errorf("load oauth repository failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "oauth repository unavailable")
	}

	client, err := repo.GetOAuthClientByClientID(req.ClientID)
	if err != nil {
		l.Errorf("load oauth client failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to load oauth client")
	}
	if client == nil {
		return nil, newOAuthError(http.StatusUnauthorized, "invalid_client", "oauth client not found")
	}
	if client.ClientSecret != "" && strings.TrimSpace(client.ClientSecret) != strings.TrimSpace(req.ClientSecret) {
		return nil, newOAuthError(http.StatusUnauthorized, "invalid_client", "client authentication failed")
	}
	if !redirectURIMatches(client, req.RedirectURI) {
		return nil, newOAuthError(http.StatusBadRequest, "invalid_request", "redirect_uri is not allowed")
	}

	now := time.Now()
	code, err := repo.ConsumeAuthorizationCode(req.Code, client.ClientID, now)
	if err != nil {
		l.Errorf("consume oauth code failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to exchange authorization code")
	}
	if code == nil {
		return nil, newOAuthError(http.StatusBadRequest, "invalid_grant", "authorization code is invalid or expired")
	}
	if !strings.EqualFold(strings.TrimSpace(code.RedirectURI), strings.TrimSpace(req.RedirectURI)) {
		return nil, newOAuthError(http.StatusBadRequest, "invalid_grant", "redirect_uri does not match authorization code")
	}

	user, err := repo.GetUserByID(code.UserID)
	if err != nil {
		l.Errorf("load oauth user failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to load oauth user")
	}
	if user == nil || !user.Enabled {
		return nil, newOAuthError(http.StatusUnauthorized, "access_denied", "oauth user is unavailable")
	}

	accessToken, err := newOpaqueToken(48)
	if err != nil {
		l.Errorf("generate oauth access token failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to generate access token")
	}

	ttl := accessTokenTTL()
	token := &relationDB.IamOAuthAccessToken{
		AccessToken: accessToken,
		ClientID:    client.ClientID,
		UserID:      user.ID,
		Scopes:      normalizeScope(firstNonEmptyString(code.Scopes, client.Scopes)),
		ExpiredAt:   now.Add(ttl),
		CreatedAt:   now,
	}
	if err := repo.CreateAccessToken(token); err != nil {
		l.Errorf("store oauth access token failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to persist access token")
	}

	return &types.OAuthTokenResp{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(ttl / time.Second),
		Scope:       token.Scopes,
	}, nil
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
