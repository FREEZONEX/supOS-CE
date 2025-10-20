package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 批量查询文件历史值
func FileHistoryBatchQueryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewFileHistoryBatchQueryLogic(r.Context(), svcCtx)
		err := l.FileHistoryBatchQuery()
		result.Http(w, r, nil, err)
	}
}
