// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package importExport

import (
	"net/http"

	"backend/internal/logic/supos/uns/importExport"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ExportHandler UNS 导出
func ExportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExportReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := importExport.NewExportLogic(r.Context(), svcCtx)
		resp, err := l.Export(w, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else if resp != nil {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
