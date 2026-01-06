package auth

import (
	"net/http"

	"backend/internal/common/constants"
	authlogic "backend/internal/logic/supos/auth"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/result"
)

// userinfo
func UserinfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie(constants.AccessTokenKey)
		var token string
		if cookie != nil {
			token = cookie.Value
		}

		l := authlogic.NewUserInfoLogic(r.Context(), svcCtx)
		resp, newCookie, err := l.UserInfo(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if newCookie != nil {
			http.SetCookie(w, newCookie)
		}
		result.Http(w, r, resp, nil)
	}
}
