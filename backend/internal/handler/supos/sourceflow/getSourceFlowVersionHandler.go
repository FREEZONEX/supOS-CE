// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sourceflow

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/result"
)

// Query current Node-RED flow version
func GetSourceFlowVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sourceflow.NewGetSourceFlowVersionLogic(r.Context(), svcCtx)
		resp, err := l.GetSourceFlowVersion()
		result.Http(w, r, resp, err)

	}
}
