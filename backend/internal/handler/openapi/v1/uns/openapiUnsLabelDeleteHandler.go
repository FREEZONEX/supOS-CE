// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"net/http"

	respx "backend/internal/httpx"
	"backend/internal/logic/openapi/v1/uns"
	"backend/internal/svc"
	"backend/internal/types"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

func OpenapiUnsLabelDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdReq
		if err := gozerohttpx.Parse(r, &req); err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
			return
		}

		l := uns.NewOpenapiUnsLabelDeleteLogic(r.Context(), svcCtx)
		resp, err := l.OpenapiUnsLabelDelete(&req)
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
