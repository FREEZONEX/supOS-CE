package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 清除所有外部topic
func ClearExternalTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewClearExternalTreeLogic(r.Context(), svcCtx)
		resp, err := l.ClearExternalTree()
		result.Http(w, r, resp, err)
	}
}
