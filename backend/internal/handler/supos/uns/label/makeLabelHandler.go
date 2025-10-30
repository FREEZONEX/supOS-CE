// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package label

import (
	"backend/internal/logic/supos/uns/label"
	"backend/internal/svc"
	"backend/internal/types"
	"encoding/json"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 文件打标签
func MakeLabelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MakeLabelReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if len(body) > 2 {
			// 手动解析 JSON Body
			if err := json.Unmarshal(body, &req.LabelList); err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		}

		l := label.NewMakeLabelLogic(r.Context(), svcCtx)
		resp, err := l.MakeLabel(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
