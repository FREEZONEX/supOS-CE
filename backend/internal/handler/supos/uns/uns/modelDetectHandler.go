package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 预先判断是否有属性关联
func ModelDetectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewModelDetectLogic(r.Context(), svcCtx)
		err := l.ModelDetect()
		result.Http(w, r, nil, err)
	}
}
