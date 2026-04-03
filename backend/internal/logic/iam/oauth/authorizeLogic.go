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

type AuthorizeLogic struct {
	baseOAuthLogic
}

type AuthorizeResult struct {
	RedirectURL string
}

func NewAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorizeLogic {
	return &AuthorizeLogic{
		baseOAuthLogic: newBaseOAuthLogic(ctx, svcCtx),
	}
}

func (l *AuthorizeLogic) Authorize(req *types.OAuthAuthorizeReq, sessionID, requestURI string) (*AuthorizeResult, error) {
	if req == nil {
		return nil, newOAuthError(http.StatusBadRequest, "invalid_request", "missing authorize request")
	}
	if !strings.EqualFold(strings.TrimSpace(req.ResponseType), "code") {
		return nil, newOAuthError(http.StatusBadRequest, "unsupported_response_type", "only authorization code flow is supported")
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
	if !redirectURIMatches(client, req.RedirectURI) {
		return nil, newOAuthError(http.StatusBadRequest, "invalid_request", "redirect_uri is not allowed")
	}

	session, err := l.getSession(sessionID)
	if err != nil {
		l.Errorf("load oauth session failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to load session")
	}
	if !l.isSessionActive(session) {
		if strings.TrimSpace(sessionID) != "" {
			l.revokeSession(sessionID)
		}
		return &AuthorizeResult{RedirectURL: loginRedirectURL(requestURI)}, nil
	}

	user, err := repo.GetUserByID(session.UserID)
	if err != nil {
		l.Errorf("load oauth user failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to load user")
	}
	if user == nil || !user.Enabled {
		l.revokeSession(sessionID)
		return &AuthorizeResult{RedirectURL: loginRedirectURL(requestURI)}, nil
	}

	codeValue, err := newOpaqueToken(32)
	if err != nil {
		l.Errorf("generate oauth code failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to generate authorization code")
	}

	now := time.Now()
	if err := repo.CreateAuthorizationCode(&relationDB.IamOAuthAuthorizationCode{
		Code:        codeValue,
		ClientID:    client.ClientID,
		UserID:      user.ID,
		RedirectURI: strings.TrimSpace(req.RedirectURI),
		Scopes:      normalizeScope(req.Scope),
		ExpiredAt:   now.Add(authorizationCodeTTL),
		CreatedAt:   now,
	}); err != nil {
		l.Errorf("create oauth authorization code failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to create authorization code")
	}
	if err := l.touchSession(sessionID); err != nil {
		l.Errorf("touch oauth session failed: %v", err)
	}

	redirectURL, err := appendAuthorizeCode(req.RedirectURI, codeValue, req.State)
	if err != nil {
		return nil, newOAuthError(http.StatusBadRequest, "invalid_request", "redirect_uri is invalid")
	}
	return &AuthorizeResult{RedirectURL: redirectURL}, nil
}
