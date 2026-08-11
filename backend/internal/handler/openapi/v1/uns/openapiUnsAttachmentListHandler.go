// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"net/http"

	respx "backend/internal/httpx"
	"backend/internal/logic/openapi/v1/uns"
	"backend/internal/svc"
	"backend/internal/types"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

func OpenapiUnsAttachmentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnsAttachmentListReq
		if err := gozerohttpx.Parse(r, &req); err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
			return
		}
		if req.Topic == "" && req.UnsId <= 0 {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "topic or unsId is required"))
			return
		}

		l := uns.NewOpenapiUnsAttachmentListLogic(r.Context(), svcCtx)
		resp, err := l.OpenapiUnsAttachmentList(&req)
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
