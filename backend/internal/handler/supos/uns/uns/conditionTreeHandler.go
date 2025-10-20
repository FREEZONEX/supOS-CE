package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 多条件分页查询树结构
func ConditionTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewConditionTreeLogic(r.Context(), svcCtx)
		err := l.ConditionTree()
		result.Http(w, r, nil, err)
	}
}
