package resource

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/resource"
	"backend/internal/svc"
)

// get
func GetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := resource.NewGetLogic(r.Context(), svcCtx)
		err := l.Get()
		result.Http(w, r, nil, err)
	}
}
