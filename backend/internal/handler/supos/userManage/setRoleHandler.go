package userManage

import (
	"backend/internal/logic/supos/userManage"
	"backend/internal/svc"
	"backend/internal/types"
	"net/http"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Assign roles to user
func SetRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserSetRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := userManage.NewSetRoleLogic(r.Context(), svcCtx)
		resp, err := l.SetRole(&req)
		result.Http(w, r, resp, err)
	}
}
