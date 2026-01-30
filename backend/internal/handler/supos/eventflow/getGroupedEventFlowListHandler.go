// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"backend/internal/logic/supos/sourceflow"
	"net/http"

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

		l := sourceflow.NewGetGroupedSourceFlowListLogic(r.Context(), svcCtx)
		req.GroupType = 2
		resp, err := l.GetGroupedSourceFlowList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
