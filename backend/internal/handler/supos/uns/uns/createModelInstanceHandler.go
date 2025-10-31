// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/share/base"
	"bytes"
	"io"
	"net/http"
	"strings"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建文件夹和文件
func CreateModelInstanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if body := r.Body; body != nil {
			bs, _ := io.ReadAll(body)
			if len(bs) > 0 {
				jsonBody := string(bs)
				jsonBody = strings.Replace(jsonBody, `"parentId":""`, `"parentId":null`, 1)
				r.Body = base.NewReadCloserWrapper(bytes.NewBuffer([]byte(jsonBody)))
			}
		}
		var req types.CreateTopicDto
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := uns.NewCreateModelInstanceLogic(r.Context(), svcCtx)
		resp, err := l.CreateModelInstance(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
