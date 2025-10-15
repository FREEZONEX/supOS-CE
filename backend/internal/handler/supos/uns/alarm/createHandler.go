package alarm

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/alarm"
	"backend/internal/svc"
)

// 创建报警规则
func CreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := alarm.NewCreateLogic(r.Context(), svcCtx)
		err := l.Create()
		result.Http(w, r, nil, err)
	}
}
