// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package system

import (
	respx "backend/internal/httpx"
	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"backend/internal/logic/core/system"
	"backend/internal/svc"
)

func ReadyzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := system.NewReadyzLogic(r.Context(), svcCtx)
		resp, err := l.Readyz()
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
