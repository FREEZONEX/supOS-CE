package example

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/example"
	"backend/internal/svc"
)

// 模拟restApi数据
func MockRestapiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := example.NewMockRestapiLogic(r.Context(), svcCtx)
		err := l.MockRestapi()
		result.Http(w, r, nil, err)
	}
}
