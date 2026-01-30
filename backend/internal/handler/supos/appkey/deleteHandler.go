package appkey

import (
	"backend/internal/logic/supos/appkey"
	"backend/internal/svc"
	"backend/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除密钥
func DeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteIDReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := appkey.NewDeleteLogic(r.Context(), svcCtx)
		var err = l.Delete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			result := &types.JsonResult{
				Code: 200,
				Msg:  "成功",
			}
			httpx.OkJsonCtx(r.Context(), w, result)
		}
	}
}
