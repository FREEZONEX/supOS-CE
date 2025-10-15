package dashboard

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
)

// bindUns
func BindUnsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dashboard.NewBindUnsLogic(r.Context(), svcCtx)
		err := l.BindUns()
		result.Http(w, r, nil, err)
	}
}
