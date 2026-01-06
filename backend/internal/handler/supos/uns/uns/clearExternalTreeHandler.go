// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"net/http"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 清除所有外部topic
func ClearExternalTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewClearExternalTreeLogic(r.Context(), svcCtx)
		resp, err := l.ClearExternalTree()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
