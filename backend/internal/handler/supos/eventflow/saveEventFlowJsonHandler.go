// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"net/http"

	"backend/internal/logic/supos/eventflow"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Persist Node-RED event flow JSON
func SaveEventFlowJsonHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EventFlowSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := eventflow.NewSaveEventFlowJsonLogic(r.Context(), svcCtx)
		err := l.SaveEventFlowJson(&req)
		result.Http(w, r, nil, err)

	}
}
