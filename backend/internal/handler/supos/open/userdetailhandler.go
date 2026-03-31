// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"net/http"

	"backend/internal/logic/supos/open"
	"backend/internal/svc"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户详情
func UserDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数中获取 username
		username := r.PathValue("username")
		if username == "" {
			httpx.ErrorCtx(r.Context(), w, errors.Parameter.WithMsg("username is required"))
			return
		}

		l := open.NewUserDetailLogic(r.Context(), svcCtx)
		resp, err := l.UserDetail(username)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
