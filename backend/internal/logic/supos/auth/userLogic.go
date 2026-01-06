package auth

import (
	"context"
	"strings"

	cache "backend/internal/common/cache"
	"backend/internal/common/vo"
	"backend/internal/svc"
)

type UserLogic struct {
	baseAuthLogic
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		baseAuthLogic: newBaseAuthLogic(ctx, svcCtx),
	}
}

func (l *UserLogic) User(sessionKey string) (any, error) {
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		if l.authDisabled() {
			return vo.Guest(), nil
		}
		return "not found user info", nil
	}

	entry, ok := l.getTokenEntry(sessionKey)
	if !ok || entry == nil || entry.Token == nil {
		if l.authDisabled() {
			return vo.Guest(), nil
		}
		l.removeSession(sessionKey)
		return "not found user info", nil
	}

	userInfo, sub, err := l.fetchUserInfo(entry.Token.AccessToken, true)
	if err != nil {
		l.Errorf("fetch user info failed: %v", err)
	}
	if userInfo == nil {
		l.removeSession(sessionKey)
		if sub != "" && cache.UserInfoCache != nil {
			cache.UserInfoCache.Delete(sub)
		}
		if l.authDisabled() {
			return vo.Guest(), nil
		}
		return "not found user info", nil
	}
	userInfo.SuperAdmin = userInfo.IsSuperAdmin()
	return userInfo, nil
}
