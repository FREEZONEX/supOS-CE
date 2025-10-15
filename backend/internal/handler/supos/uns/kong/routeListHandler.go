package kong

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/kong"
	"backend/internal/svc"
)

// 获取路由
func RouteListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := kong.NewRouteListLogic(r.Context(), svcCtx)
		err := l.RouteList()
		result.Http(w, r, nil, err)
	}
}
