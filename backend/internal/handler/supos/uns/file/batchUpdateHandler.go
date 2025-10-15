package file

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/file"
	"backend/internal/svc"
)

// 批量修改文件值
func BatchUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewBatchUpdateLogic(r.Context(), svcCtx)
		err := l.BatchUpdate()
		result.Http(w, r, nil, err)
	}
}
