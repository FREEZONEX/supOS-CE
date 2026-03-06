package attachment

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"backend/internal/logic/supos/attachment"
	"backend/internal/svc"
	"backend/internal/types"
)

// DownloadAttachmentHandler 下载附件
func DownloadAttachmentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadAttachmentRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := attachment.NewDownloadAttachmentLogic(r.Context(), svcCtx)
		// 下载文件直接写入response
		if err := l.DownloadAttachment(&req, w); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		// 成功时不需要返回JSON，直接返回文件内容
	}
}