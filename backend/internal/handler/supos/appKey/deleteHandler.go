package appKey

import (
	"net/http"
	"strconv"

	"backend/internal/logic/supos/appkey"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
)

// 删除密钥
func DeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("id无效"))
			return
		}

		l := appKey.NewDeleteLogic(r.Context(), svcCtx)
		err = l.Delete(id)
		result.Http(w, r, nil, err)
	}
}
