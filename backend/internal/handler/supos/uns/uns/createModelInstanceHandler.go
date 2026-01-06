// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
	"backend/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建文件夹和文件
func CreateModelInstanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTopicDto
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := uns.NewCreateModelInstanceLogic(r.Context(), svcCtx)
		resp, err := l.CreateModelInstance(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
