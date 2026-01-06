package auth

import (
	"net/http"

	"backend/internal/common/constants"
	authlogic "backend/internal/logic/supos/auth"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/result"
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
		result.Http(w, r, nil, err)
	}
}
