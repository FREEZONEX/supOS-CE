package oauth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"backend/internal/svc"
	"backend/internal/types"
)

const (
	defaultPortainerClientID      = "portainer"
	defaultPortainerBootstrapUser = "tier0"
)

type UserInfoLogic struct {
	baseOAuthLogic
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		baseOAuthLogic: newBaseOAuthLogic(ctx, svcCtx),
	}
}

func oauthEnvOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value != "" {
		return value
	}
	return fallback
}

func isPortainerClient(clientID string) bool {
	return strings.EqualFold(strings.TrimSpace(clientID), oauthEnvOrDefault("IAM_PORTAINER_CLIENT_ID", defaultPortainerClientID))
}

func (l *UserInfoLogic) UserInfo(accessToken string) (*types.OAuthUserInfoResp, error) {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return nil, newOAuthError(http.StatusUnauthorized, "invalid_token", "access token is required")
	}

	repo, err := l.repo()
	if err != nil || repo == nil {
		l.Errorf("load oauth repository failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "oauth repository unavailable")
	}

	token, err := repo.GetAccessToken(accessToken)
	if err != nil {
		l.Errorf("load oauth access token failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to load access token")
	}
	if token == nil || token.RevokedAt != nil || !token.ExpiredAt.After(time.Now()) {
		return nil, newOAuthError(http.StatusUnauthorized, "invalid_token", "access token is invalid or expired")
	}

	user, err := repo.GetUserByID(token.UserID)
	if err != nil {
		l.Errorf("load oauth user failed: %v", err)
		return nil, newOAuthError(http.StatusInternalServerError, "server_error", "failed to load oauth user")
	}
	if user == nil || !user.Enabled {
		return nil, newOAuthError(http.StatusUnauthorized, "invalid_token", "oauth user is unavailable")
	}

	name := strings.TrimSpace(user.DisplayName)
	if name == "" {
		name = strings.TrimSpace(user.Username)
	}

	if isPortainerClient(token.ClientID) {
		username := oauthEnvOrDefault("IAM_BOOTSTRAP_USERNAME", defaultPortainerBootstrapUser)
		return &types.OAuthUserInfoResp{
			Sub:               username,
			PreferredUsername: username,
			Email:             strings.TrimSpace(oauthEnvOrDefault("IAM_BOOTSTRAP_EMAIL", "")),
			Name:              username,
		}, nil
	}

	return &types.OAuthUserInfoResp{
		Sub:               user.ID,
		PreferredUsername: strings.TrimSpace(user.Username),
		Email:             strings.TrimSpace(user.Email),
		Name:              name,
	}, nil
}
