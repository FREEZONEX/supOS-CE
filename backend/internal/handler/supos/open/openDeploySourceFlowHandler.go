package open

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 部署 Flow
func OpenDeploySourceFlowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenSourceFlowDeployReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		resp, err := sourceflow.NewDeploySourceFlowLogic(r.Context(), svcCtx).DeploySourceFlow(&types.SourceFlowDeployReq{
			ID:    req.ID,
			Flows: req.Flows,
		})
		result.Http(w, r, resp, err)
	}
}
