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

func FlowDataSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FlowDataReq
		if err := gozerohttpx.Parse(r, &req); err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
			return
		}

		l := flow.NewFlowDataSaveLogic(r.Context(), svcCtx)
		resp, err := l.FlowDataSave(&req)
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
