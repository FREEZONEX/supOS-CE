package auth

import (
	"context"
	"net/http"
	"strings"

	cache "backend/internal/common/cache"
	"backend/internal/common/constants"
	"backend/internal/common/vo"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/errors"
)

type UserInfoLogic struct {
	baseAuthLogic
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		baseAuthLogic: newBaseAuthLogic(ctx, svcCtx),
	}
}

func (l *UserInfoLogic) UserInfo(sessionKey string) (*vo.UserInfoVo, *http.Cookie, error) {
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		return nil, nil, errors.NotLogin.WithMsg("cookie missing")
	}

	entry, ok := l.getTokenEntry(sessionKey)
	if !ok || entry == nil || entry.Token == nil {
		l.removeSession(sessionKey)
		return nil, nil, errors.NotLogin.WithMsg("token not found")
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
		return nil, nil, errors.NotLogin.WithMsg("not found user info")
	}

	if cache.TokenCache != nil {
		cache.TokenCache.Refresh(sessionKey)
	}
	cookie := &http.Cookie{
		Name:     constants.AccessTokenKey,
		Value:    sessionKey,
		Path:     "/",
		MaxAge:   constants.CookieMaxAge,
		HttpOnly: false,
		Secure:   false,
	}
	return userInfo, cookie, nil
}
