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

// 模板实例附件删除
func AttachmentDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AttachmentDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := attachment.NewAttachmentDeleteLogic(r.Context(), svcCtx)
		err := l.AttachmentDelete(&req)
		result.Http(w, r, nil, err)
	}
}
