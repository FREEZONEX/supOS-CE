// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"net/http"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建文件夹和文件
func ModelCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnsCreateTopicDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		//l := uns.NewModelCreateLogic(r.Context(), svcCtx)
		//resp, err := l.ModelCreate(&req)
		//if err != nil {
		httpx.ErrorCtx(r.Context(), w, nil)
		//} else {
		//	httpx.OkJsonCtx(r.Context(), w, resp)
		//}
	}
}
