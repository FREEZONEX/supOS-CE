package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 文件或文件夹修改订阅
func ModelSubscribeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewModelSubscribeLogic(r.Context(), svcCtx)
		err := l.ModelSubscribe()
		result.Http(w, r, nil, err)
	}
}
