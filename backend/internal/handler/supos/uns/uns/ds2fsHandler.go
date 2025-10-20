package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 外部数据源表的字段定义转uns字段定义
func Ds2fsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewDs2fsLogic(r.Context(), svcCtx)
		err := l.Ds2fs()
		result.Http(w, r, nil, err)
	}
}
