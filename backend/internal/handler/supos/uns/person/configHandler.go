package person

import (
	"net/http"

	"backend/internal/logic/supos/uns/person"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取个人配置
func ConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPersonConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := person.NewConfigLogic(r.Context(), svcCtx)
		resp, err := l.Config(&req)
		result.Http(w, r, resp, err)
	}
}
