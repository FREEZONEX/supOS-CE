// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package topology

import (
	"net/http"

	"backend/internal/logic/supos/uns/topology"
	"backend/internal/svc"
	"backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取实例拓扑状态
func GetInstanceTopologyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetInstanceTopologyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := topology.NewGetInstanceTopologyLogic(r.Context(), svcCtx)
		resp, err := l.GetInstanceTopology(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
