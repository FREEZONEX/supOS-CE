// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package system

import (
	"net/http"

	"backend/internal/logic/supos/uns/system"
	"backend/internal/svc"
	"backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取UNS定义
func GetUnsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDefRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := system.NewGetUnsLogic(r.Context(), svcCtx)
		resp, err := l.GetUns(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
