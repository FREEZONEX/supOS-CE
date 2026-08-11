// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package info

import (
	"net/http"

	respx "backend/internal/httpx"
	"backend/internal/logic/openapi/v1/info"
	"backend/internal/svc"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

func OpenapiInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := info.NewOpenapiInfoLogic(r.Context(), svcCtx)
		resp, err := l.OpenapiInfo()
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
