package global

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/global"
	"backend/internal/svc"
)

// 全局数据导入
func DataImportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := global.NewDataImportLogic(r.Context(), svcCtx)
		err := l.DataImport()
		result.Http(w, r, nil, err)
	}
}
