package example

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/example"
	"backend/internal/svc"
)

// 模拟订单数据
func MockRestapiOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := example.NewMockRestapiOrderLogic(r.Context(), svcCtx)
		err := l.MockRestapiOrder()
		result.Http(w, r, nil, err)
	}
}
