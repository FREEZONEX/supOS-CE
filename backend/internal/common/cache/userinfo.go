package cache

import (
	"backend/internal/common/constants"
	"backend/internal/common/vo"
	"time"
)

func newUserInfoCache() (*ManagedCache[*vo.UserInfoVo], error) {
	ttl := time.Duration(constants.TokenMaxAge) * time.Second
	return NewManagedCache[*vo.UserInfoVo](defaultUserInfoCacheCapacity, ttl)
}
