package dashboard

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
)

// 取消置顶
func UnmarkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dashboard.NewUnmarkLogic(r.Context(), svcCtx)
		err := l.Unmark()
		result.Http(w, r, nil, err)
	}
}
