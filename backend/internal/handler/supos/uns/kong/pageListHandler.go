package kong

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/kong"
	"backend/internal/svc"
)

// 查询报警列表
func PageListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := kong.NewPageListLogic(r.Context(), svcCtx)
		err := l.PageList()
		result.Http(w, r, nil, err)
	}
}
