package dashboard

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
)

// isExist
func IsExistHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dashboard.NewIsExistLogic(r.Context(), svcCtx)
		err := l.IsExist()
		result.Http(w, r, nil, err)
	}
}
