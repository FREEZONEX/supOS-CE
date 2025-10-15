package mount

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/mount"
	"backend/internal/svc"
)

// 获取挂载数据源
func SourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := mount.NewSourceLogic(r.Context(), svcCtx)
		err := l.Source()
		result.Http(w, r, nil, err)
	}
}
