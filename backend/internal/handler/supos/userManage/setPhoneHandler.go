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

// Update phone number
func SetPhoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PhoneUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := userManage.NewSetPhoneLogic(r.Context(), svcCtx)
		resp, err := l.SetPhone(&req)
		result.Http(w, r, resp, err)
	}
}
