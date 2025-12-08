// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"io"
	"net/http"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 外部JSON定义转树结构uns字段定义
func ParseJson2TreeUnsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := uns.NewParseJson2TreeUnsLogic(r.Context(), svcCtx)
		resp, err := l.ParseJson2TreeUns(body)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
