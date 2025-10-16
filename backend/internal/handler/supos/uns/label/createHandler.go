package label

import (
	"backend/internal/logic/supos/uns/label"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 创建标签
func CreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnsLabelCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := label.NewCreateLogic(r.Context(), svcCtx)
		resp, err := l.Create(&req)
		result.Http(w, r, resp, err)
	}
}
