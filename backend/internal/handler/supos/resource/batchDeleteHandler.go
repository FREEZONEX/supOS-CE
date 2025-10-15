package resource

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/resource"
	"backend/internal/svc"
)

// batch
func BatchDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := resource.NewBatchDeleteLogic(r.Context(), svcCtx)
		err := l.BatchDelete()
		result.Http(w, r, nil, err)
	}
}
