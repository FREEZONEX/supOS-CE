// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"net/http"
	"strconv"

	"backend/internal/logic/supos/open"
	"backend/internal/svc"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除模板
func DeleteTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数中获取 id
		idStr := r.PathValue("id")
		if idStr == "" {
			httpx.ErrorCtx(r.Context(), w, errors.Parameter.WithMsg("id is required"))
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.Parameter.WithMsg("invalid id format"))
			return
		}

		l := open.NewDeleteTemplateLogic(r.Context(), svcCtx)
		resp, err := l.DeleteTemplate(id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
