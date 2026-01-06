package resource

import (
	"backend/internal/logic/supos/resource"
	"backend/internal/svc"
	"backend/internal/types"
	"net/http"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Batch update resources
func BatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []types.BatchUpdateResource
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := resource.NewBatchLogic(r.Context(), svcCtx)
		err := l.Batch(&req)
		result.Http(w, r, nil, err)
	}
}
