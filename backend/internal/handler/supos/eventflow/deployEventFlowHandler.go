// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"net/http"

	"backend/internal/logic/supos/eventflow"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Deploy an event flow
func DeployEventFlowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EventFlowDeployReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := eventflow.NewDeployEventFlowLogic(r.Context(), svcCtx)
		resp, err := l.DeployEventFlow(&req)
		result.Http(w, r, resp, err)

	}
}
