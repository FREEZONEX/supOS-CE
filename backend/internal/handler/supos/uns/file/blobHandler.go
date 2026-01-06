package file

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/file"
	"backend/internal/svc"
)

// 获取文件BLOB类型的值
func BlobHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewBlobLogic(r.Context(), svcCtx)
		err := l.Blob()
		result.Http(w, r, nil, err)
	}
}
