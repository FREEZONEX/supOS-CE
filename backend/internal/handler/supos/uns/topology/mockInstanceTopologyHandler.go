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

// 设置实例拓扑状态(mock)
func MockInstanceTopologyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MockInstanceTopologyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := topology.NewMockInstanceTopologyLogic(r.Context(), svcCtx)
		err := l.MockInstanceTopology(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
