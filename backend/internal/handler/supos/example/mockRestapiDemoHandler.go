package example

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/example"
	"backend/internal/svc"
)

// 模拟电器厂IT数据
func MockRestapiDemoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := example.NewMockRestapiDemoLogic(r.Context(), svcCtx)
		err := l.MockRestapiDemo()
		result.Http(w, r, nil, err)
	}
}
