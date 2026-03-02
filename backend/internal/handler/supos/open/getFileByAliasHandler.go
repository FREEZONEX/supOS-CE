// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"net/http"

	"backend/internal/logic/supos/open"
	"backend/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 别名查询文件详情
func GetFileByAliasHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := open.NewGetFileByAliasLogic(r.Context(), svcCtx)
		resp, err := l.GetFileByAlias()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
