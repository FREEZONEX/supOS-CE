package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 删除指定路径下的所有文件夹和文件
func DeletePathHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewDeletePathLogic(r.Context(), svcCtx)
		err := l.DeletePath()
		result.Http(w, r, nil, err)
	}
}
