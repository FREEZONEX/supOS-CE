// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package service_api

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow/service_api"
	"backend/internal/svc"
	"backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Batch query flows by UNS aliases
func BatchSourceFlowByAliasesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SourceFlowBatchAliasReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := service_api.NewBatchSourceFlowByAliasesLogic(r.Context(), svcCtx)
		resp, err := l.BatchSourceFlowByAliases(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
