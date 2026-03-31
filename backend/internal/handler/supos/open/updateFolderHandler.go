// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"net/http"

	"backend/internal/logic/supos/open"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 修改文件夹
func UpdateFolderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateOpenApiFolderDto
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 从路径参数中获取 alias
		alias := r.PathValue("alias")
		if alias == "" {
			httpx.ErrorCtx(r.Context(), w, errors.Parameter.WithMsg("alias is required"))
			return
		}

		l := open.NewUpdateFolderLogic(r.Context(), svcCtx)
		resp, err := l.UpdateFolder(alias, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
