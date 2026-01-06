// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"net/http"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 批量创建文件夹和文件(node-red导入专用)
func CreateModelsForNodeRedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []*types.CreateUnsNodeRedDto
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := uns.NewCreateModelsForNodeRedLogic(r.Context(), svcCtx)
		resp, err := l.CreateModelsForNodeRed(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
