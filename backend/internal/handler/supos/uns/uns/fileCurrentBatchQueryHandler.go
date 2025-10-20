package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 批量查询文件实时值
func FileCurrentBatchQueryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewFileCurrentBatchQueryLogic(r.Context(), svcCtx)
		err := l.FileCurrentBatchQuery()
		result.Http(w, r, nil, err)
	}
}
