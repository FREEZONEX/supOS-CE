// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package importExport

import (
	"backend/internal/logic/supos/auth"
	"backend/internal/logic/supos/uns/importExport"
	"backend/internal/svc"
	"backend/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ExportPathHandler UNS 导出,返回文件路径
func ExportPathHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExportReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		req.UserId = auth.ResolveUserID(r.Context())
		req.Language = types.GetAcceptLanguage(r)
		l := importExport.NewExportPathLogic(r.Context(), svcCtx)
		resp, err := l.ExportPath(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
