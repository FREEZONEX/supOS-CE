package cache

import (
	"time"

	"backend/internal/common/constants"
	authdto "backend/internal/common/dto/auth"
)

// TokenCacheEntry keeps the raw token payload alongside useful metadata.
type TokenCacheEntry struct {
	Token    *authdto.AccessTokenDto
	Raw      map[string]any
	CachedAt time.Time
}

func newTokenCache() (*ManagedCache[*TokenCacheEntry], error) {
	ttl := time.Duration(constants.TokenMaxAge) * time.Second
	return NewManagedCache[*TokenCacheEntry](defaultTokenCacheCapacity, ttl)
}
