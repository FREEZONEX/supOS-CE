package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 根据别名集合批量删除文件夹和文件
func BatchAliasHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewBatchAliasLogic(r.Context(), svcCtx)
		err := l.BatchAlias()
		result.Http(w, r, nil, err)
	}
}
