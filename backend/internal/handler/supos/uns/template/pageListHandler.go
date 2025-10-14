package template

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/template"
	"backend/internal/svc"
)

// 查询模板列表
func PageListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := template.NewPageListLogic(r.Context(), svcCtx)
		err := l.PageList()
		result.Http(w, r, nil, err)
	}
}
