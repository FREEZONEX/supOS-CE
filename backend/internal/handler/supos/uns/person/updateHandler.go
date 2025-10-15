package person

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/person"
	"backend/internal/svc"
)

// 设置个人配置
func UpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := person.NewUpdateLogic(r.Context(), svcCtx)
		err := l.Update()
		result.Http(w, r, nil, err)
	}
}
