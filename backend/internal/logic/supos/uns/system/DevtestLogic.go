package system

import (
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/utils/dbpool"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"encoding/json"
	"strconv"
	"unicode"

	"github.com/zeromicro/go-zero/core/logx"
)

func Devtest(ctx context.Context, params map[string][]string) (resp map[string]interface{}) {
	resp = map[string]interface{}{}
	logx.Debug("Devtest: Debug")
	logx.Info("Devtest: Info")

	user := apiutil.GetUserFromContext(ctx)
	userBs, _ := json.Marshal(user)
	logx.WithContext(ctx).Info("WithCtxLog: 当前用户:", string(userBs))

	resp["user"] = user
	resp["poolStats"] = dbpool.Stats()

	if len(params) > 0 {
		defService := spring.GetBean[serviceApi.IUnsDefinitionService]()
		if aliasList := params["uns"]; len(aliasList) > 0 {
			for _, a := range aliasList {
				if len(a) > 0 {
					var uns *types.CreateTopicDto
					if unicode.IsDigit(rune(a[0])) {
						id, er := strconv.ParseInt(a, 10, 64)
						if er == nil {
							uns = defService.GetDefinitionById(id)
						}
					} else {
						uns = defService.GetDefinitionByAlias(a)
					}
					if uns != nil {
						resp[a] = uns
					}
				}
			}
		}
	}
	return
}
