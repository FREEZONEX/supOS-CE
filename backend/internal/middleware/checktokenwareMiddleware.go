package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"backend/internal/common/constants"
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/vo"
	iamrepo "backend/internal/repo/iam"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type CheckTokenWareMiddleware struct {
	defaultHome string
}

func NewCheckTokenWareMiddleware(defaultHome string) *CheckTokenWareMiddleware {
	return &CheckTokenWareMiddleware{
		defaultHome: defaultHome,
	}
}

func (m *CheckTokenWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authDisabled := func() bool {
			flag := os.Getenv("SYS_OS_AUTH_ENABLE")
			return strings.EqualFold(flag, "false")
		}()
		ctx := context.WithValue(r.Context(), "lang", r.Header.Get("Accept-Language"))
		cookieToken, err := apiutil.GetCookie(r, constants.AccessTokenKey)
		if err != nil || cookieToken == "" {
			if authDisabled {
				r = r.WithContext(apiutil.SetUserInContext(ctx, vo.Guest()))
				next(w, r)
				return
			}
			httpx.WriteJson(w, http.StatusUnauthorized, result.Error(errors.NotLogin.Code, "not logged in"))
			return
		}

		repo, repoErr := iamrepo.NewAuthRepo(ctx)
		if repoErr != nil || repo == nil {
			httpx.WriteJson(w, http.StatusUnauthorized, result.Error(errors.NotLogin.Code, "session store unavailable"))
			return
		}

		session, sessionErr := repo.GetSession(cookieToken)
		if sessionErr == nil && session != nil && session.RevokedAt == nil && session.ExpiredAt.After(time.Now()) {
			user, userErr := repo.BuildUserInfo(ctx, session.UserID, m.defaultHome)
			if userErr == nil && user != nil {
				now := time.Now()
				_ = repo.TouchSession(cookieToken, now, now.Add(time.Duration(constants.TokenMaxAge)*time.Second))
				r = r.WithContext(apiutil.SetUserInContext(ctx, user))
				next(w, r)
				return
			}
		}

		if session != nil {
			_ = repo.RevokeSession(cookieToken, time.Now())
		}
		if authDisabled {
			r = r.WithContext(apiutil.SetUserInContext(ctx, vo.Guest()))
			next(w, r)
			return
		}
		httpx.WriteJson(w, http.StatusUnauthorized, result.Error(errors.NotLogin.Code, "token expired"))
	}
}
