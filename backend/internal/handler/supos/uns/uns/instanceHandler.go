package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 查询文件详情
func InstanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewInstanceLogic(r.Context(), svcCtx)
		err := l.Instance()
		result.Http(w, r, nil, err)
	}
}
