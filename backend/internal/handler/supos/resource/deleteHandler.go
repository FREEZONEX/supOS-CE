package resource

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/resource"
	"backend/internal/svc"
)

// delete
func DeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := resource.NewDeleteLogic(r.Context(), svcCtx)
		err := l.Delete()
		result.Http(w, r, nil, err)
	}
}
