// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package nodered

import (
	"net/http"

	"backend/internal/logic/supos/nodered"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 代理 NodeRed /flows 接口
func ProxyNodeRedFlowsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookieID := ""
		if c, err := r.Cookie("flowId"); err == nil {
			cookieID = c.Value
		}

		l := nodered.NewProxyNodeRedFlowsLogic(r.Context(), svcCtx)
		resp, err := l.ProxyNodeRedFlows(cookieID)
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
