package resource

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/resource"
	"backend/internal/svc"
)

// batch
func BatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := resource.NewBatchLogic(r.Context(), svcCtx)
		err := l.Batch()
		result.Http(w, r, nil, err)
	}
}
