package oauth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"backend/internal/svc"
)

type LogoutLogic struct {
	baseOAuthLogic
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		baseOAuthLogic: newBaseOAuthLogic(ctx, svcCtx),
	}
}

func (l *LogoutLogic) Logout(accessToken string) error {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return newOAuthError(http.StatusBadRequest, "invalid_request", "access token is required")
	}

	repo, err := l.repo()
	if err != nil || repo == nil {
		l.Errorf("load oauth repository failed: %v", err)
		return newOAuthError(http.StatusInternalServerError, "server_error", "oauth repository unavailable")
	}
	if err := repo.RevokeAccessToken(accessToken, time.Now()); err != nil {
		l.Errorf("revoke oauth access token failed: %v", err)
		return newOAuthError(http.StatusInternalServerError, "server_error", "failed to revoke access token")
	}
	return nil
}
