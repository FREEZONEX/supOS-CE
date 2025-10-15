package dashboard

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
)

// edit
func EditHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dashboard.NewEditLogic(r.Context(), svcCtx)
		err := l.Edit()
		result.Http(w, r, nil, err)
	}
}
