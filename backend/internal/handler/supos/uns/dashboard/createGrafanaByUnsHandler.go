package dashboard

import (
	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
	"backend/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// createGrafanaByUns
func CreateGrafanaByUnsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AliasPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dashboard.NewCreateGrafanaByUnsLogic(r.Context(), svcCtx)
		resp, err := l.CreateGrafanaByUns(req.Alias)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
