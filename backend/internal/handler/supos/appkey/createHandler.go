package appkey

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/appkey"
	"backend/internal/svc"
)

// 创建密钥
func CreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := appkey.NewCreateLogic(r.Context(), svcCtx)
		err := l.Create()
		result.Http(w, r, nil, err)
	}
}
