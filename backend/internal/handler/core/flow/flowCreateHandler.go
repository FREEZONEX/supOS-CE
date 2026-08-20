// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	respx "backend/internal/httpx"
	"backend/internal/logic/core/flow"
	"backend/internal/svc"
	"backend/internal/types"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func FlowCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FlowSaveReq
		if err := gozerohttpx.Parse(r, &req); err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
			return
		}

		l := flow.NewFlowCreateLogic(r.Context(), svcCtx)
		resp, err := l.FlowCreate(&req)
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
