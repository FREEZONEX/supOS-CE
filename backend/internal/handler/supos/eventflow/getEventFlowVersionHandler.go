// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"net/http"

	"backend/internal/logic/supos/eventflow"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/result"
)

// Query current event flow version
func GetEventFlowVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := eventflow.NewGetEventFlowVersionLogic(r.Context(), svcCtx)
		resp, err := l.GetEventFlowVersion()
		result.Http(w, r, resp, err)

	}
}
