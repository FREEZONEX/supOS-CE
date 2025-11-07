package i18n

import (
	"backend/internal/logic/supos/i18n"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 返回支持的语言
func GetLanguagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLanguagesReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := i18n.NewGetLanguagesLogic(r.Context(), svcCtx)
		resp, err := l.GetLanguages(&req)
		result.Http(w, r, resp, err)
	}
}
