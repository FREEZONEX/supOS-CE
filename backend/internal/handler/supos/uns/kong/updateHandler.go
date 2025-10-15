package kong

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/kong"
	"backend/internal/svc"
)

// 更新报警规则
func UpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := kong.NewUpdateLogic(r.Context(), svcCtx)
		err := l.Update()
		result.Http(w, r, nil, err)
	}
}
