package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 删除前预先判断是否有被引用对象
func DetectIfRemoveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewDetectIfRemoveLogic(r.Context(), svcCtx)
		err := l.DetectIfRemove()
		result.Http(w, r, nil, err)
	}
}
