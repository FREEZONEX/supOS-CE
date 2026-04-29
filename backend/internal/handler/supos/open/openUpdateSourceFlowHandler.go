package open

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新 Flow
func OpenUpdateSourceFlowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenSourceFlowUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		err := sourceflow.NewUpdateSourceFlowLogic(r.Context(), svcCtx).UpdateSourceFlow(&types.SourceFlowUpdateReq{
			ID:          req.ID,
			FlowName:    req.FlowName,
			Description: req.Description,
		})
		result.Http(w, r, nil, err)
	}
}
