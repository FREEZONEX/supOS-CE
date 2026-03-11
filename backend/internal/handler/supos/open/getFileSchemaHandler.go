// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"net/http"

	"backend/internal/logic/supos/open"
	"backend/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询文件schema 元数据结构
func GetFileSchemaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := open.NewGetFileSchemaLogic(r.Context(), svcCtx)
		resp, err := l.GetFileSchema()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
