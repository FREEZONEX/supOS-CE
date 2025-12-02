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

// UploadHandler UNS 上传文件
func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := types.FormFile(r, "file")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := importExport.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
