package devtools

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/devtools"
	"backend/internal/svc"
)

// logs
func LogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := devtools.NewLogsLogic(r.Context(), svcCtx)
		err := l.Logs()
		result.Http(w, r, nil, err)
	}
}
