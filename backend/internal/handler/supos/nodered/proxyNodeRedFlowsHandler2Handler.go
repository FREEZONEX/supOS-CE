// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package nodered

import (
	"net/http"

	"backend/internal/logic/supos/nodered"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 代理 NodeRed /flows 接口（备用路径）
func ProxyNodeRedFlowsHandler2Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := nodered.NewProxyNodeRedFlowsHandler2Logic(r.Context(), svcCtx)
		resp, err := l.ProxyNodeRedFlowsHandler2()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
