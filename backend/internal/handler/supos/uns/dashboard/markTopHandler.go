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

// 置顶 Dashboard
func MarkTopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarkTopRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dashboard.NewMarkTopLogic(r.Context(), svcCtx)
		resp, err := l.MarkTop(req.ID, getUserID(r))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
