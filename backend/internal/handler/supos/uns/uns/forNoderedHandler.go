package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 批量创建文件夹和文件(node-red导入专用)
func ForNoderedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewForNoderedLogic(r.Context(), svcCtx)
		err := l.ForNodered()
		result.Http(w, r, nil, err)
	}
}
