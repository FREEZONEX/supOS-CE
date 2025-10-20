package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 批量写文件实时值
func FileCurrentBatchUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewFileCurrentBatchUpdateLogic(r.Context(), svcCtx)
		err := l.FileCurrentBatchUpdate()
		result.Http(w, r, nil, err)
	}
}
