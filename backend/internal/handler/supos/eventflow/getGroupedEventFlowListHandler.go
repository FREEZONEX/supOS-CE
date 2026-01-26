// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"net/http"

	"backend/internal/logic/supos/eventflow"
	"backend/internal/svc"
	"backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 分页按分组获取event flow列表
func GetGroupedEventFlowListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupPageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := eventflow.NewGetGroupedEventFlowListLogic(r.Context(), svcCtx)
		resp, err := l.GetGroupedEventFlowList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
