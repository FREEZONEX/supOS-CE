// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"net/http"

	"backend/internal/logic/supos/open"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 模板实例附件上传
func AttachmentUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const maxBodySize = 10 << 20
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		if err := r.ParseMultipartForm(maxBodySize); err != nil { // 10MB max memory
			httpx.ErrorCtx(r.Context(), w, errors.Parameter.WithMsg("uns.importDocumentMax"))
			return
		}

		var req types.AttachmentUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := open.NewAttachmentUploadLogic(r.Context(), svcCtx, r)
		resp, err := l.AttachmentUpload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			result.Http(w, r, resp, err)
		}
	}
}
