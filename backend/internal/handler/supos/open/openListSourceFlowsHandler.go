package open

import (
	"net/http"

	"backend/internal/logic/supos/sourceflow"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询 Flow 列表
func OpenListSourceFlowsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenSourceFlowListQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		resp, err := sourceflow.NewListSourceFlowsLogic(r.Context(), svcCtx).ListSourceFlows(&types.SourceFlowListQuery{
			Keyword:   req.Keyword,
			OrderCode: req.OrderCode,
			IsAsc:     req.IsAsc,
			PageNo:    req.PageNo,
			PageSize:  req.PageSize,
		})
		result.HttpWithoutWrap(w, r, resp, err)
	}
}
