// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package i18n

import (
	"net/http"

	respx "backend/internal/httpx"
	"backend/internal/logic/core/i18n"
	"backend/internal/svc"
	"backend/internal/types"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

// 获取系统语言包
func I18nMessagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.I18nMessagesReq
		if err := gozerohttpx.Parse(r, &req); err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "入参不正确:"+err.Error()))
			return
		}

		l := i18n.NewI18nMessagesLogic(r.Context(), svcCtx)
		resp, err := l.I18nMessages(&req)
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
