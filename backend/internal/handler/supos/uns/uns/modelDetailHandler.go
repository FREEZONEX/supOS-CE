package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 查询文件夹详情
func ModelDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewModelDetailLogic(r.Context(), svcCtx)
		err := l.ModelDetail()
		result.Http(w, r, nil, err)
	}
}
