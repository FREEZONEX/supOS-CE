package external

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/external"
	"backend/internal/svc"
)

// 搜索外部topic主题树，默认整个树
func TreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := external.NewTreeLogic(r.Context(), svcCtx)
		err := l.Tree()
		result.Http(w, r, nil, err)
	}
}
