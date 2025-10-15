package dashboard

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
)

// 创建
func CreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dashboard.NewCreateLogic(r.Context(), svcCtx)
		err := l.Create()
		result.Http(w, r, nil, err)
	}
}
