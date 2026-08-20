// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package auth

import (
	respx "backend/internal/httpx"
	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"backend/internal/logic/core/auth"
	"backend/internal/svc"
)

func CurrentUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewCurrentUserLogic(r.Context(), svcCtx)
		resp, err := l.CurrentUser()
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
