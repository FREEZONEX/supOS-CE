package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 清除所有外部topic
func ExternalTopicClearHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewExternalTopicClearLogic(r.Context(), svcCtx)
		err := l.ExternalTopicClear()
		result.Http(w, r, nil, err)
	}
}
