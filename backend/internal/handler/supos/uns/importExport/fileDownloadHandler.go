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

// FileDownloadHandler 文件下载
func FileDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileDownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := importExport.NewFileDownloadLogic(r.Context(), svcCtx)
		err := l.FileDownload(&req, r, w)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
