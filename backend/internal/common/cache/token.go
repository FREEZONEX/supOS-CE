package cache

import (
	"time"

	"backend/internal/common/constants"
	"backend/share/clients"
)

// TokenCacheEntry keeps the raw token payload fetched from Keycloak alongside useful metadata.
type TokenCacheEntry struct {
	Token    *clients.AccessTokenDto
	Raw      map[string]any
	CachedAt time.Time
}

func newTokenCache() (*ManagedCache[*TokenCacheEntry], error) {
	ttl := time.Duration(constants.TokenMaxAge) * time.Second
	return NewManagedCache[*TokenCacheEntry](defaultTokenCacheCapacity, ttl)
}
