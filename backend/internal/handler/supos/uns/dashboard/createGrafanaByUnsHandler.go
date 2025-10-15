package dashboard

import (
	"backend/internal/logic/supos/uns/dashboard"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// createGrafanaByUns
func CreateGrafanaByUnsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateGrafanaByUnsReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := dashboard.NewCreateGrafanaByUnsLogic(r.Context(), svcCtx)
		err := l.CreateGrafanaByUns(&req)
		result.Http(w, r, nil, err)
	}
}
