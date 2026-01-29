package appKey

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/appkey"
	"backend/internal/svc"
)

// 查询密钥列表
func ListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := appKey.NewListLogic(r.Context(), svcCtx)
		resp, err := l.List()
		result.Http(w, r, resp, err)
	}
}
