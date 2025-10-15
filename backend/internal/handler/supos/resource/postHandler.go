package resource

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/resource"
	"backend/internal/svc"
)

// post
func PostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := resource.NewPostLogic(r.Context(), svcCtx)
		err := l.Post()
		result.Http(w, r, nil, err)
	}
}
