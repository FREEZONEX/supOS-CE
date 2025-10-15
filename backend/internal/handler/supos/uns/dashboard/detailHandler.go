package dashboard

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
)

// 获取详情
func DetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dashboard.NewDetailLogic(r.Context(), svcCtx)
		err := l.Detail()
		result.Http(w, r, nil, err)
	}
}
