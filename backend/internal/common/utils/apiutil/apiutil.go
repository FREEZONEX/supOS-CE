package apiutil

import (
	"backend/internal/common/vo"
	"context"
	"net/http"
	"os"
)

// userContextKey is an unexported type to prevent context key collisions.
type userContextKey string

const (
	// UserKey is the key for storing user info in the context.
	UserKey = userContextKey("user")
)

// SetUserInContext returns a new request with the user information stored in its context.
// This is the idiomatic Go way to pass values through context, typically used in middleware.
// It does NOT modify the original request.
func SetUserInContext(r *http.Request, user *vo.UserInfoVo) *http.Request {
	ctx := context.WithValue(r.Context(), UserKey, user)
	return r.WithContext(ctx)
}

// GetUserFromContext retrieves user information from the request's context.
// If auth is disabled and no user is in the context, it returns a default "guest" user.
func GetUserFromContext(r *http.Request) *vo.UserInfoVo {
	// Try to get user from context
	if user, ok := r.Context().Value(UserKey).(*vo.UserInfoVo); ok && user != nil {
		return user
	}

	// If auth is disabled, mock a guest user
	authEnable := os.Getenv("SYS_OS_AUTH_ENABLE")
	if authEnable == "false" || authEnable == "" {
		return vo.Guest()
	}

	return nil
}

// GetCookie retrieves a cookie from the request.
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
