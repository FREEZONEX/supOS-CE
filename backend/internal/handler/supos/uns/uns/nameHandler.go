package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 修改文件夹或文件名称
func NameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewNameLogic(r.Context(), svcCtx)
		err := l.Name()
		result.Http(w, r, nil, err)
	}
}
