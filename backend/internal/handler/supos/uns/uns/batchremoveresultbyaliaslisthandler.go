package uns

import (
	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 根据别名集合批量删除文件夹和文件
func BatchRemoveResultByAliasListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchRemoveUnsDto
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := uns.NewBatchRemoveResultByAliasListLogic(r.Context(), svcCtx)
		resp, err := l.BatchRemoveResultByAliasList(&req)
		result.Http(w, r, resp, err)
	}
}
