package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns"
	"backend/internal/svc"
)

// 获取系统配置
func SystemConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewSystemConfigLogic(r.Context(), svcCtx)
		err := l.SystemConfig()
		result.Http(w, r, nil, err)
	}
}
