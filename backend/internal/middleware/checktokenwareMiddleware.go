package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"backend/internal/common/authsvc"
	cache "backend/internal/common/cache"
	"backend/internal/common/constants"
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/vo"
	"backend/share/clients"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type CheckTokenWareMiddleware struct {
	keycloak    *clients.KeycloakClient
	defaultHome string
	realm       string
}

func NewCheckTokenWareMiddleware(kc *clients.KeycloakClient, defaultHome, realm string) *CheckTokenWareMiddleware {
	flag := os.Getenv("SYS_OS_AUTH_ENABLE")
	fmt.Println("SYS_OS_AUTH_ENABLE:", flag)
	return &CheckTokenWareMiddleware{
		keycloak:    kc,
		defaultHome: defaultHome,
		realm:       realm,
	}
}

func (m *CheckTokenWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authDisabled := func() bool {
			flag := os.Getenv("SYS_OS_AUTH_ENABLE")
			return flag == "" || strings.EqualFold(flag, "false")
		}()

		cookieToken, err := apiutil.GetCookie(r, constants.AccessTokenKey)
		if err != nil || cookieToken == "" {
			if authDisabled {
				r = r.WithContext(apiutil.SetUserInContext(r.Context(), vo.Guest()))
				next(w, r)
				return
			}
			httpx.WriteJson(w, http.StatusUnauthorized, result.Error(errors.NotLogin.Code, "not logged in"))
			return
		}

		if cache.TokenCache == nil {
			httpx.WriteJson(w, http.StatusUnauthorized, result.Error(errors.NotLogin.Code, "token cache not initialized"))
			return
		}

		entry, ok := cache.TokenCache.Get(cookieToken)
		if !ok || entry == nil || entry.Token == nil || entry.Token.AccessToken == "" {
			if authDisabled {
				r = r.WithContext(apiutil.SetUserInContext(r.Context(), vo.Guest()))
				next(w, r)
				return
			}
			httpx.WriteJson(w, http.StatusUnauthorized, result.Error(errors.NotLogin.Code, "token expired"))
			return
		}

		cache.TokenCache.Refresh(cookieToken)

		user, _, fetchErr := authsvc.FetchUserInfo(r.Context(), m.keycloak, entry.Token.AccessToken, true, m.defaultHome, m.realm)
		if fetchErr != nil {
			if authDisabled {
				r = r.WithContext(apiutil.SetUserInContext(r.Context(), vo.Guest()))
				next(w, r)
				return
			}
			httpx.WriteJson(w, http.StatusUnauthorized, result.Error(errors.NotLogin.Code, fetchErr.Error()))
			return
		}

		if user == nil {
			if authDisabled {
				r = r.WithContext(apiutil.SetUserInContext(r.Context(), vo.Guest()))
				next(w, r)
				return
			}
			httpx.WriteJson(w, http.StatusUnauthorized, result.Error(errors.NotLogin.Code, "user not found"))
			return
		}

		r = r.WithContext(apiutil.SetUserInContext(r.Context(), user))
		next(w, r)
	}
}
