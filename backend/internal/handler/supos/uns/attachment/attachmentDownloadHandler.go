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

// 模板实例附件下载
func AttachmentDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AttachmentDownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := attachment.NewAttachmentDownloadLogic(r.Context(), svcCtx)
		err:=l.AttachmentDownload(&req,w,r)
		if err != nil {
			result.Http(w, r, nil, err)
			return
		}
	}
}
