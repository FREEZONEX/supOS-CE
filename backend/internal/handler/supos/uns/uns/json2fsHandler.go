package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 外部JSON定义转uns字段定义
func Json2fsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewJson2fsLogic(r.Context(), svcCtx)
		err := l.Json2fs()
		result.Http(w, r, nil, err)
	}
}
