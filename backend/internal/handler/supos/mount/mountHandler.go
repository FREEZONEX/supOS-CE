package mount

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/mount"
	"backend/internal/svc"
)

// 手动挂载
func MountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := mount.NewMountLogic(r.Context(), svcCtx)
		err := l.Mount()
		result.Http(w, r, nil, err)
	}
}
