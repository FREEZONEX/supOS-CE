package config

import "time"

// TokenCacheConfig represents token cache configuration.
type TokenCacheConfig struct {
	TokenCacheTTL    time.Duration // Token cache TTL (默认1小时 = 60*60*1000ms)
	UserInfoCacheTTL time.Duration // User info cache TTL (默认无限期 = Long.MAX_VALUE)
}

// NewTokenCacheConfig returns default token cache configuration
func NewTokenCacheConfig() *TokenCacheConfig {
	return &TokenCacheConfig{
		TokenCacheTTL:    time.Hour,        // 1 hour = 60*60*1000ms
		UserInfoCacheTTL: time.Duration(0), // No expiration = Long.MAX_VALUE
	}
}

// GetTokenCacheTTL returns token cache TTL in milliseconds
func (t *TokenCacheConfig) GetTokenCacheTTL() int64 {
	return int64(t.TokenCacheTTL / time.Millisecond)
}

// GetUserInfoCacheTTL returns user info cache TTL in milliseconds
func (t *TokenCacheConfig) GetUserInfoCacheTTL() int64 {
	if t.UserInfoCacheTTL == 0 {
		return int64(^uint64(0) >> 1) // Go equivalent of Long.MAX_VALUE
	}
	return int64(t.UserInfoCacheTTL / time.Millisecond)
}
