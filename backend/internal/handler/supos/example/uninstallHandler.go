package example

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/example"
	"backend/internal/svc"
)

// 卸载
func UninstallHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := example.NewUninstallLogic(r.Context(), svcCtx)
		err := l.Uninstall()
		result.Http(w, r, nil, err)
	}
}
