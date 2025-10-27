package auth

import (
	"net/http"

	"backend/internal/common/constants"
	authlogic "backend/internal/logic/supos/auth"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// logout
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie(constants.AccessTokenKey)
		var token string
		if cookie != nil {
			token = cookie.Value
		}
		l := authlogic.NewLogoutLogic(r.Context(), svcCtx)
		err := l.Logout(token)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
