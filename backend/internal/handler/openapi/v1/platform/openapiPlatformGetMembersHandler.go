// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package platform

import (
	"net/http"

	respx "backend/internal/httpx"
	"backend/internal/logic/openapi/v1/platform"
	"backend/internal/svc"
	"backend/internal/types"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

func OpenapiPlatformGetMembersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PlatformGetMembersReq
		if err := gozerohttpx.Parse(r, &req); err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
			return
		}

		l := platform.NewOpenapiPlatformGetMembersLogic(r.Context(), svcCtx)
		resp, err := l.OpenapiPlatformGetMembers(&req)
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
