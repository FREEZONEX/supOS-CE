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
		l := nodered.NewProxyNodeRedFlowsLogic(r.Context(), svcCtx)
		resp, err := l.ProxyNodeRedFlows()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
