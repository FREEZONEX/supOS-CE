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

// 修改模板订阅
func UpdateSubscribeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTemplateSubscribeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := template.NewUpdateSubscribeLogic(r.Context(), svcCtx)
		resp, err := l.UpdateSubscribe(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
