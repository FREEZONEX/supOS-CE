package i18n

import (
	"backend/internal/logic/supos/uns/i18n"
	"backend/internal/svc"
	"backend/internal/types"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// 获取i18n
func ReadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUnsI18nMessagesReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := i18n.NewReadFileLogic(r.Context(), svcCtx)
		resp, err := l.ReadFile(&req)
		result.Http(w, r, resp, err)
	}
}
