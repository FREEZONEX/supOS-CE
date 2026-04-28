package open

import (
	"encoding/json"
	"net/http"

	"backend/internal/logic/supos/sourceflow/service_api"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询 Flow Data
func OpenGetSourceFlowDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenSourceFlowDataQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		raw, err := service_api.NewProxySourceFlowsLogic(r.Context(), svcCtx).ProxySourceFlows(req.ID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		var resp types.OpenSourceFlowDataResp
		if err := json.Unmarshal([]byte(raw), &resp); err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.System.WithMsg("error.sys.systemError").AddDetail(err.Error()))
			return
		}

		result.HttpWithoutWrap(w, r, resp, nil)
	}
}
