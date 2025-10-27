package cache

import (
	"backend/internal/common/vo"
	"time"

	"github.com/maypok86/otter/v2"
)

const (
	defaultTokenCacheCapacity    = 1024
	defaultUserInfoCacheCapacity = 1024
)

var (
	TokenCache    *ManagedCache[*TokenCacheEntry]
	UserInfoCache *ManagedCache[*vo.UserInfoVo]
)

// InitAuthCaches prepares the token and user info caches with the provided TTL settings.
func InitCaches() (err error) {

	TokenCache, err = newTokenCache()
	if err != nil {
		return err
	}

	UserInfoCache, err = newUserInfoCache()
	if err != nil {
		return err
	}

	return nil
}

// ManagedCache is a thin wrapper around otter.Cache that provides a consistent API for different value types.
type ManagedCache[T any] struct {
	cache *otter.Cache[string, T]
}

// NewManagedCache builds a cache with the provided capacity and TTL.
func NewManagedCache[T any](capacity int, ttl time.Duration) (*ManagedCache[T], error) {
	opts := &otter.Options[string, T]{
		MaximumSize: capacity,
	}
	if ttl > 0 {
		opts.ExpiryCalculator = otter.ExpiryWriting[string, T](ttl)
	}
	cache, err := otter.New(opts)
	if err != nil {
		return nil, err
	}
	return &ManagedCache[T]{cache: cache}, nil
}

// Set stores a value by key.
func (m *ManagedCache[T]) Set(key string, value T) {
	if m == nil || m.cache == nil {
		return
	}
	m.cache.Set(key, value)
}

// Get retrieves a cached value, if present.
func (m *ManagedCache[T]) Get(key string) (T, bool) {
	if m == nil || m.cache == nil {
		var zero T
		return zero, false
	}
	return m.cache.GetIfPresent(key)
}

// Refresh re-applies the TTL for an entry by writing the current value back.
func (m *ManagedCache[T]) Refresh(key string) {
	if m == nil || m.cache == nil {
		return
	}
	if value, ok := m.cache.GetIfPresent(key); ok {
		m.cache.Set(key, value)
	}
}

// Delete removes a cached entry.
func (m *ManagedCache[T]) Delete(key string) {
	if m == nil || m.cache == nil {
		return
	}
	m.cache.Invalidate(key)
}
