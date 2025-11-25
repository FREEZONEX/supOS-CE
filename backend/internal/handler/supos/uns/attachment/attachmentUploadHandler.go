package attachment

import (
	"backend/internal/logic/supos/uns/attachment"
	"backend/internal/svc"
	"backend/internal/types"
	"net/http"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 模板实例附件上传
func AttachmentUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AttachmentUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := attachment.NewAttachmentUploadLogic(r.Context(), svcCtx, r)
		resp, err := l.AttachmentUpload(&req)
		result.Http(w, r, resp, err)
	}
}
