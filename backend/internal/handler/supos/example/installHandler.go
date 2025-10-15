package example

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/example"
	"backend/internal/svc"
)

// 安装接口
func InstallHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := example.NewInstallLogic(r.Context(), svcCtx)
		err := l.Install()
		result.Http(w, r, nil, err)
	}
}
