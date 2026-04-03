package auth

import (
	"net/http"

	authlogic "backend/internal/logic/supos/auth"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// login
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, err)
			return
		}

		l := authlogic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			result.Http(w, r, nil, err)
			return
		}
		if resp != nil && resp.Cookie != nil {
			http.SetCookie(w, resp.Cookie)
		}
		result.Http(w, r, resp.User, nil)
	}
}
