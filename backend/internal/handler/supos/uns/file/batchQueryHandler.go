package file

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/file"
	"backend/internal/svc"
)

// 批量查询文件实时值
func BatchQueryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewBatchQueryLogic(r.Context(), svcCtx)
		err := l.BatchQuery()
		result.Http(w, r, nil, err)
	}
}
