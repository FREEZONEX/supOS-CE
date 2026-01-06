// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package label

import (
	"net/http"

	"backend/internal/logic/supos/uns/label"
	"backend/internal/svc"
	"backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 文件打单个标签
func MakeSingleLabelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MakeSingleLabelReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := label.NewMakeSingleLabelLogic(r.Context(), svcCtx)
		resp, err := l.MakeSingleLabel(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
