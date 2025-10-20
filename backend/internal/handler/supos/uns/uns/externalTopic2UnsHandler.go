package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 外部topic转UNS
func ExternalTopic2UnsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewExternalTopic2UnsLogic(r.Context(), svcCtx)
		err := l.ExternalTopic2Uns()
		result.Http(w, r, nil, err)
	}
}
