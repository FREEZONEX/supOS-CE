package uns

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"backend/internal/logic/supos/uns/uns"
	"backend/internal/svc"
)

// 校验指定文件夹夹是否已存在文件夹、文件名称
func NameDuplicationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := uns.NewNameDuplicationLogic(r.Context(), svcCtx)
		err := l.NameDuplication()
		result.Http(w, r, nil, err)
	}
}
