package dashboard

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
)

// getByUns
func GetByUnsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dashboard.NewGetByUnsLogic(r.Context(), svcCtx)
		err := l.GetByUns()
		result.Http(w, r, nil, err)
	}
}
