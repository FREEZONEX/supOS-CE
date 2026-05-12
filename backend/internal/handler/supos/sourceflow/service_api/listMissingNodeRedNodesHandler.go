// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package service_api

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow/service_api"
	"backend/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// List missing Node-RED nodes across all source flow tabs
func ListMissingNodeRedNodesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := service_api.NewListMissingNodeRedNodesLogic(r.Context(), svcCtx)
		resp, err := l.ListMissingNodeRedNodes(r.URL.Query().Get("flowType"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
