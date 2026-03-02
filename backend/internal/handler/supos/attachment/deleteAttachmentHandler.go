// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package attachment

import (
	"net/http"

	"backend/internal/logic/supos/attachment"
	"backend/internal/svc"
	"backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Delete attachment
func DeleteAttachmentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteAttachmentRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := attachment.NewDeleteAttachmentLogic(r.Context(), svcCtx)
		resp, err := l.DeleteAttachment(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
