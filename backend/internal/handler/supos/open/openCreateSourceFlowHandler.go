package open

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建 Flow
func OpenCreateSourceFlowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenSourceFlowCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		resp, err := sourceflow.NewCreateSourceFlowLogic(r.Context(), svcCtx).CreateSourceFlow(&types.SourceFlowCreateReq{
			FlowName:    req.FlowName,
			Description: req.Description,
			Template:    req.Template,
			GroupId:     req.GroupId,
		})
		result.Http(w, r, resp, err)
	}
}
