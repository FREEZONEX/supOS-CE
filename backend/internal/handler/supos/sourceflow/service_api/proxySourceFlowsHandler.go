// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package service_api

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow/service_api"
	"backend/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Proxy Node-RED /flows endpoint using cookie scoped id
func ProxySourceFlowsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookieID := ""
		if c, err := r.Cookie("sup_flow_id"); err == nil {
			cookieID = c.Value
		}

		l := service_api.NewProxySourceFlowsLogic(r.Context(), svcCtx)
		resp, err := l.ProxySourceFlows(cookieID)
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
