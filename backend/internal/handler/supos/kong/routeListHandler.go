// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package kong

import (
	"net/http"

	"backend/internal/logic/supos/kong"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取简化的路由列表
func RouteListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := kong.NewRouteListLogic(r.Context(), svcCtx)
		resp, err := l.RouteList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
