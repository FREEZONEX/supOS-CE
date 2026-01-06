package global

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/global"
	"backend/internal/svc"
)

// 确认导出记录
func UserExportRecordConfirmHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := global.NewUserExportRecordConfirmLogic(r.Context(), svcCtx)
		err := l.UserExportRecordConfirm()
		result.Http(w, r, nil, err)
	}
}
