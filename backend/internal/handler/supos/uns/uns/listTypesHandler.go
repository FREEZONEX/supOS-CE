// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"net/http"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 枚举数据类型
func ListTypesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewListTypesLogic(r.Context(), svcCtx)
		resp, err := l.ListTypes()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
