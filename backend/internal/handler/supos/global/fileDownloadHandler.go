package global

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/global"
	"backend/internal/svc"
)

// 根据路径下载文件
func FileDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := global.NewFileDownloadLogic(r.Context(), svcCtx)
		err := l.FileDownload()
		result.Http(w, r, nil, err)
	}
}
