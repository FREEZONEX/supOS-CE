package userManage

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/userManage"
	"backend/internal/svc"
)

// List available roles
func RoleListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := userManage.NewRoleListLogic(r.Context(), svcCtx)
		resp, err := l.RoleList()
		result.Http(w, r, resp, err)
	}
}
