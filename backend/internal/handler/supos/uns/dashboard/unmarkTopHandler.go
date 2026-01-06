// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"net/http"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 取消置顶 Dashboard
func UnmarkTopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnmarkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dashboard.NewUnmarkTopLogic(r.Context(), svcCtx)
		resp, err := l.UnmarkTop(req.ID, getUserID(r))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
