// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"net/http"

	"backend/internal/logic/supos/open"
	"backend/internal/svc"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询模板详情
func TemplateDetailByAliasHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数中获取 alias
		alias := r.PathValue("alias")
		if alias == "" {
			httpx.ErrorCtx(r.Context(), w, errors.Parameter.WithMsg("alias is required"))
			return
		}

		l := open.NewTemplateDetailByAliasLogic(r.Context(), svcCtx)
		resp, err := l.TemplateDetailByAlias(alias)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
