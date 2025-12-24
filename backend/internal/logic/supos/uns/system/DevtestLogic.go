package system

import (
	"backend/internal/common/utils/apiutil"
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

func Devtest(ctx context.Context, params map[string][]string) (resp map[string]interface{}) {
	resp = map[string]interface{}{}
	logx.Debug(ctx, "Devtest: Debug")
	logx.Info(ctx, "Devtest: Info")

	user := apiutil.GetUserFromContext(ctx)
	userBs, _ := json.Marshal(user)
	logx.WithContext(ctx).Info("WithCtxLog: 当前用户:", string(userBs))

	resp["user"] = user
	return
}
