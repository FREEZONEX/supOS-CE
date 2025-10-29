package dashboard

import (
	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
	"net/http"

	"gitee.com/unitedrhino/share/result"
)

// createGrafanaByUns
func CreateGrafanaByUnsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dashboard.NewCreateGrafanaByUnsLogic(r.Context(), svcCtx)
		err := l.CreateGrafanaByUns()
		result.Http(w, r, nil, err)
	}
}
