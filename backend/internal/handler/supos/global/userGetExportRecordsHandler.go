package global

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/global"
	"backend/internal/svc"
)

// 获取导出记录
func UserGetExportRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := global.NewUserGetExportRecordsLogic(r.Context(), svcCtx)
		err := l.UserGetExportRecords()
		result.Http(w, r, nil, err)
	}
}
