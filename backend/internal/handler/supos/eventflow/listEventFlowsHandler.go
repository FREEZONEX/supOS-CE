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

// List event flows with optional fuzzy search
func ListEventFlowsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EventFlowListQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := eventflow.NewListEventFlowsLogic(r.Context(), svcCtx)
		resp, err := l.ListEventFlows(&req)
		result.HttpWithoutWrap(w, r, resp, err)

	}
}
