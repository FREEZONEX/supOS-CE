// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Create a mock flow from UNS path
func CreateSourceMockFlowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SourceFlowMockReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sourceflow.NewCreateSourceMockFlowLogic(r.Context(), svcCtx)
		resp, err := l.CreateSourceMockFlow(&req)
		result.Http(w, r, resp, err)

	}
}
