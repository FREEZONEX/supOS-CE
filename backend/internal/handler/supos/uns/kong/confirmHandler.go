package kong

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/kong"
	"backend/internal/svc"
)

// 确认报警
func ConfirmHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := kong.NewConfirmLogic(r.Context(), svcCtx)
		err := l.Confirm()
		result.Http(w, r, nil, err)
	}
}
