// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package auth

import (
	"net/http"

	respx "backend/internal/httpx"
	"backend/internal/logic/openapi/v1/auth"
	"backend/internal/svc"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

func OpenapiWhoamiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewOpenapiWhoamiLogic(r.Context(), svcCtx)
		resp, err := l.OpenapiWhoami()
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
