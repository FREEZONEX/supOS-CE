package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/internal/contextx"
	"backend/internal/repo"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	iam       *repo.IAMRepo
	apiKeys   *repo.APIKeyRepo
	jwtSecret []byte
}

type LoginResult struct {
	Token               string
	User                repo.User
	ForceChangePassword bool
	ResourceKeys        []string
	UIActions           []string
}

func New(ctx context.Context, jwtSecret string) *Service {
	return &Service{
		iam:       repo.NewIAMRepo(ctx),
		apiKeys:   repo.NewAPIKeyRepo(ctx),
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (LoginResult, error) {
	user, err := s.iam.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return LoginResult{}, errUnauthorized
		}
		return LoginResult{}, err
	}
	if user.Status != 1 {
		return LoginResult{}, errUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return LoginResult{}, errUnauthorized
	}
	if err := s.iam.MarkLogin(ctx, user.ID); err != nil {
		return LoginResult{}, err
	}
	keys, err := s.iam.ResourceKeysForUser(ctx, user.ID)
	if err != nil {
		return LoginResult{}, err
	}
	uiActions, err := s.iam.UIActionsForResourceKeys(ctx, keys)
	if err != nil {
		return LoginResult{}, err
	}
	token, err := s.issueToken(user)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		Token:               token,
		User:                user,
		ForceChangePassword: UsesDefaultPassword(user),
		ResourceKeys:        keys,
		UIActions:           uiActions,
	}, nil
}

func (s *Service) CurrentUser(ctx context.Context, userID int64) (repo.User, []string, []string, bool, error) {
	user, err := s.iam.GetUserByID(ctx, userID)
	if err != nil {
		return repo.User{}, nil, nil, false, err
	}
	if user.Status != 1 {
		return repo.User{}, nil, nil, false, errUnauthorized
	}
	keys, err := s.iam.ResourceKeysForUser(ctx, user.ID)
	if err != nil {
		return repo.User{}, nil, nil, false, err
	}
	uiActions, err := s.iam.UIActionsForResourceKeys(ctx, keys)
	if err != nil {
		return repo.User{}, nil, nil, false, err
	}
	return user, keys, uiActions, UsesDefaultPassword(user), nil
}

func (s *Service) AuthenticateRequest(r *http.Request) (contextx.Subject, error) {
	token := bearerToken(r)
	if token == "" {
		if cookie, err := r.Cookie("edge_session"); err == nil {
			token = cookie.Value
		}
	}
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		return contextx.Subject{}, errUnauthorized
	}
	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return contextx.Subject{}, errUnauthorized
	}
	idFloat, ok := claims["uid"].(float64)
	if !ok {
		return contextx.Subject{}, errUnauthorized
	}
	user, keys, uiActions, _, err := s.CurrentUser(r.Context(), int64(idFloat))
	if err != nil {
		return contextx.Subject{}, errUnauthorized
	}
	return contextx.Subject{
		UserID: user.ID, UserName: user.UserName, Email: user.Email, AuthType: "session",
		ResourceKeys: toSet(keys), UIActions: uiActions,
	}, nil
}

func (s *Service) AuthenticateAPIKey(ctx context.Context, raw string) (contextx.Subject, error) {
	if raw == "" {
		return contextx.Subject{}, errUnauthorized
	}
	key, err := s.apiKeys.ValidateAPIKeyDetail(ctx, raw)
	if err != nil {
		return contextx.Subject{}, errUnauthorized
	}
	user, ownerKeys, uiActions, _, err := s.CurrentUser(ctx, key.OwnerID)
	if err != nil {
		return contextx.Subject{}, errUnauthorized
	}
	effective := intersect(ownerKeys, key.ResourceKeys)
	return contextx.Subject{
		UserID:       user.ID,
		UserName:     user.UserName,
		Email:        user.Email,
		AuthType:     "apiKey",
		APIKeyID:     key.ID,
		APIKeyName:   key.Name,
		APIKeyPrefix: key.KeyPrefix,
		KeyType:      key.KeyType,
		ResourceKeys: toSet(effective),
		UIActions:    uiActions,
	}, nil
}

// RevalidateMQTTAPIKey reloads an established MQTT session's API key by its
// non-secret record ID. Invalid or disabled credentials are returned as
// valid=false, while repository failures remain errors so callers do not
// disconnect healthy clients during a transient database outage.
func (s *Service) RevalidateMQTTAPIKey(ctx context.Context, keyID int64) (contextx.Subject, bool, error) {
	key, err := s.apiKeys.ValidateAPIKeyByIDDetail(ctx, keyID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return contextx.Subject{}, false, nil
		}
		return contextx.Subject{}, false, err
	}
	user, err := s.iam.GetUserByID(ctx, key.OwnerID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return contextx.Subject{}, false, nil
		}
		return contextx.Subject{}, false, err
	}
	if user.Status != 1 {
		return contextx.Subject{}, false, nil
	}
	return contextx.Subject{
		UserID:       user.ID,
		UserName:     user.UserName,
		Email:        user.Email,
		AuthType:     "apiKey",
		APIKeyID:     key.ID,
		APIKeyName:   key.Name,
		APIKeyPrefix: key.KeyPrefix,
		KeyType:      key.KeyType,
		ResourceKeys: toSet(key.ResourceKeys),
	}, true, nil
}

func UsesDefaultPassword(user repo.User) bool {
	if strings.TrimSpace(user.Password) == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("tier0")) == nil
}

func (s *Service) issueToken(user repo.User) (string, error) {
	claims := jwt.MapClaims{
		"uid": user.ID,
		"sub": strconv.FormatInt(user.ID, 10),
		"exp": time.Now().Add(12 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func bearerToken(r *http.Request) string {
	value := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(value), "bearer ") {
		return strings.TrimSpace(value[7:])
	}
	return ""
}

func APIKeyFromRequest(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("X-API-Key")); value != "" {
		return value
	}
	if value := strings.TrimSpace(r.URL.Query().Get("apikey")); value != "" {
		return value
	}
	if value := strings.TrimSpace(r.URL.Query().Get("apiKey")); value != "" {
		return value
	}
	value := strings.TrimSpace(r.Header.Get("Authorization"))
	lower := strings.ToLower(value)
	if strings.HasPrefix(lower, "apikey ") {
		return strings.TrimSpace(value[7:])
	}
	return bearerToken(r)
}

func toSet(keys []string) map[string]struct{} {
	out := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		out[key] = struct{}{}
	}
	return out
}

func intersect(left, right []string) []string {
	rightSet := toSet(right)
	var out []string
	for _, key := range left {
		if _, ok := rightSet[key]; ok {
			out = append(out, key)
		}
	}
	return out
}

var errUnauthorized = errors.New("unauthorized")
