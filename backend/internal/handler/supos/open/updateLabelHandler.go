// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"net/http"
	"strconv"

	"backend/internal/logic/supos/open"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 修改标签
func UpdateLabelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateLabelDto
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

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

		l := open.NewUpdateLabelLogic(r.Context(), svcCtx)
		resp, err := l.UpdateLabel(id, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
