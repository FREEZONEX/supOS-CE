package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 修改文件夹或文件明细
func ModelUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewModelUpdateLogic(r.Context(), svcCtx)
		err := l.ModelUpdate()
		result.Http(w, r, nil, err)
	}
}
