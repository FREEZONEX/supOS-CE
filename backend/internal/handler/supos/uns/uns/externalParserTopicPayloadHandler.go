package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 外部topic payload解析
func ExternalParserTopicPayloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewExternalParserTopicPayloadLogic(r.Context(), svcCtx)
		err := l.ExternalParserTopicPayload()
		result.Http(w, r, nil, err)
	}
}
