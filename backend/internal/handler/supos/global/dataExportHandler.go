package global

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/global"
	"backend/internal/svc"
)

// 全局数据导出
func DataExportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := global.NewDataExportLogic(r.Context(), svcCtx)
		err := l.DataExport()
		result.Http(w, r, nil, err)
	}
}
