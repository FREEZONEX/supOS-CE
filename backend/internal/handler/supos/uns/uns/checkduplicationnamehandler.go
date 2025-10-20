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

// 校验指定文件夹夹是否已存在文件夹、文件名称
func CheckDuplicationNameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckDuplicationNameReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := uns.NewCheckDuplicationNameLogic(r.Context(), svcCtx)
		resp, err := l.CheckDuplicationName(&req)
		result.Http(w, r, resp, err)
	}
}
