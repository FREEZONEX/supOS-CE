package auth

import (
	"net/http"

	"backend/internal/common/constants"
	authlogic "backend/internal/logic/supos/auth"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/result"
)

// user
func UserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie(constants.AccessTokenKey)
		var token string
		if cookie != nil {
			token = cookie.Value
		}
		l := authlogic.NewUserLogic(r.Context(), svcCtx)
		resp, err := l.User(token)
		result.Http(w, r, resp, err)
		// if err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// } else {
		// 	httpx.OkJsonCtx(r.Context(), w, resp)
		// }

	}
}
