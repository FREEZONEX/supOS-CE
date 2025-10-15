package person

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/person"
	"backend/internal/svc"
)

// 获取个人配置
func ConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := person.NewConfigLogic(r.Context(), svcCtx)
		err := l.Config()
		result.Http(w, r, nil, err)
	}
}
