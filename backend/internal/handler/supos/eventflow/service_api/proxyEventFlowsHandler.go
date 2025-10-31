// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package service_api

import (
	"net/http"

	"backend/internal/logic/supos/eventflow"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Proxy Node-RED /flows endpoint using cookie scoped id
func ProxyEventFlowsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookieID := ""
		if c, err := r.Cookie("sup_event_flow_id"); err == nil {
			cookieID = c.Value
		}

		l := eventflow.NewProxyEventFlowsLogic(r.Context(), svcCtx)
		resp, err := l.ProxyEventFlows(cookieID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if resp == "" {
			resp = "[]"
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(resp))
	}
}
