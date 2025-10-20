package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 枚举数据类型
func TypesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewTypesLogic(r.Context(), svcCtx)
		err := l.Types()
		result.Http(w, r, nil, err)
	}
}
