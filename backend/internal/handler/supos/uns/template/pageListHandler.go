// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package template

import (
	"net/http"

	"backend/internal/logic/supos/uns/template"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询模板列表
func PageListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TemplateQueryVo
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := template.NewPageListLogic(r.Context(), svcCtx)
		resp, err := l.PageList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
