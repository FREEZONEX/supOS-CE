package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 枚举数据类型
func ListTypesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewListTypesLogic(r.Context(), svcCtx)
		resp, err := l.ListTypes()
		result.Http(w, r, resp, err)
	}
}
