// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package label

import (
	"backend/share/base"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"unicode"

	"backend/internal/logic/supos/uns/label"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 文件取消标签
func CancelLabelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CancelLabelReq
		// 解析查询参数
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 手动解析 JSON Body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if len(body) > 2 {
			if body[1] == '"' {
				var labelIds []string
				if err := json.Unmarshal(body, &labelIds); err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
				req.LabelIds = base.FilterAndMap[string, int64](labelIds, func(e string) (int64, bool) {
					num, er := strconv.ParseInt(e, 10, 64)
					return num, er == nil
				})
			} else if unicode.IsDigit(rune(body[1])) {
				if err := json.Unmarshal(body, &req.LabelIds); err != nil {
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
			}
		}
		l := label.NewCancelLabelLogic(r.Context(), svcCtx)
		resp, err := l.CancelLabel(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
