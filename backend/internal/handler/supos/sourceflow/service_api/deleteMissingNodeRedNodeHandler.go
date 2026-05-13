// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package service_api

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow/service_api"
	"backend/internal/svc"
	"backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Delete one missing Node-RED node by id and location
func DeleteMissingNodeRedNodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NodeRedMissingNodeDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if req.FlowType == "" {
			req.FlowType = r.URL.Query().Get("flowType")
		}

		l := service_api.NewDeleteMissingNodeRedNodeLogic(r.Context(), svcCtx)
		resp, err := l.DeleteMissingNodeRedNode(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
