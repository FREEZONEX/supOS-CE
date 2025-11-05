package dashboard

import (
	"net/http"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/logic/supos/uns/dashboard/model"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// edit
func EditHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.DashboardModel
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dashboard.NewEditLogic(r.Context(), svcCtx)
		err := l.Edit(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
