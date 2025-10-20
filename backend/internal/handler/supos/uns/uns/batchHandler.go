package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 批量创建文件夹和文件
func BatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewBatchLogic(r.Context(), svcCtx)
		err := l.Batch()
		result.Http(w, r, nil, err)
	}
}
