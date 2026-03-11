// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"net/http"

	"backend/internal/logic/supos/app"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// List all installed applications
func ListInstalledAppsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := app.NewListInstalledAppsLogic(r.Context(), svcCtx)
		resp, err := l.ListInstalledApps()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
