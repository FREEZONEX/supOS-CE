// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"backend/internal/logic/supos/open"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func GetMembersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMembersReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("parameter.invalid:"+err.Error()))
			return
		}

		l := open.NewGetMembersLogic(r.Context(), svcCtx)
		resp, err := l.GetMembers(&req)
		result.Http(w, r, resp, err)
	}
}
