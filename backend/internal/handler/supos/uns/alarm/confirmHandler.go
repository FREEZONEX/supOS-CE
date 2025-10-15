package alarm

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/alarm"
	"backend/internal/svc"
)

// 确认报警
func ConfirmHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := alarm.NewConfirmLogic(r.Context(), svcCtx)
		err := l.Confirm()
		result.Http(w, r, nil, err)
	}
}
