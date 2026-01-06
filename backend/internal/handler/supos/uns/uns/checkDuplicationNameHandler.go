// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"net/http"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 校验指定文件夹夹是否已存在文件夹、文件名称
func CheckDuplicationNameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckDuplicationNameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := uns.NewCheckDuplicationNameLogic(r.Context(), svcCtx)
		resp, err := l.CheckDuplicationName(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
