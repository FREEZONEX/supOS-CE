package example

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/example"
	"backend/internal/svc"
)

// 初始化发电机数据
func MockDemoInitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := example.NewMockDemoInitLogic(r.Context(), svcCtx)
		err := l.MockDemoInit()
		result.Http(w, r, nil, err)
	}
}
