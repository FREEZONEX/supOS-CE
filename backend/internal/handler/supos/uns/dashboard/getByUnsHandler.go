package dashboard

import (
	"net/http"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// getByUns
func GetByUnsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetByUnsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dashboard.NewGetByUnsLogic(r.Context(), svcCtx)
		resp, err := l.GetByUns(req.UnsAlias)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
