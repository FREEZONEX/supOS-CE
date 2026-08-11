// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	respx "backend/internal/httpx"
	"backend/internal/logic/openapi/v1/uns"
	"backend/internal/svc"
	"backend/internal/types"
	"net/http"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

func OpenapiUnsBindFlowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnsBindFlowReq
		if err := gozerohttpx.Parse(r, &req); err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
			return
		}
		if req.UnsId <= 0 || req.FlowId <= 0 {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "unsId and flowId are required"))
			return
		}

		l := uns.NewOpenapiUnsBindFlowLogic(r.Context(), svcCtx)
		resp, err := l.OpenapiUnsBindFlow(&req)
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
