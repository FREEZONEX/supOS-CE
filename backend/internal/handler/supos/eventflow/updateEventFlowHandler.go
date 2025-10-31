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

// Update event flow metadata
func UpdateEventFlowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EventFlowUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := eventflow.NewUpdateEventFlowLogic(r.Context(), svcCtx)
		err := l.UpdateEventFlow(&req)
		result.Http(w, r, nil, err)

	}
}
