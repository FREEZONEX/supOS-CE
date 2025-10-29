package system

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/system"
	"backend/internal/svc"
)

// 获取系统配置
func SystemConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := system.NewSystemConfigLogic(r.Context(), svcCtx)
		resp, err := l.SystemConfig()
		result.Http(w, r, resp, err)
	}
}
