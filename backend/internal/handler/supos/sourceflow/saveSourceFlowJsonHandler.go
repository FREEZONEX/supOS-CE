// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Persist Node-RED flow JSON
func SaveSourceFlowJsonHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SourceFlowSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sourceflow.NewSaveSourceFlowJsonLogic(r.Context(), svcCtx)
		err := l.SaveSourceFlowJson(&req)
		result.Http(w, r, nil, err)

	}
}
