package attachment

import (
	"backend/internal/logic/supos/uns/attachment"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 获取模板实例附件列表
func ListAttachmentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListAttachmentReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := attachment.NewListAttachmentLogic(r.Context(), svcCtx)
		resp, err := l.ListAttachment(&req)
		result.Http(w, r, resp, err)
	}
}
