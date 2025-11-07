package I18nUtils

import (
	"context"

	"gitee.com/unitedrhino/share/i18ns"
)

func GetMessage(k string, args ...any) string {
	return i18ns.LocalizeMsg(k, args)
}

func GetMessageWithCtx(ctx context.Context, k string, args ...any) string {
	return i18ns.LocalizeMsgWithCtx(ctx, k, args)
}
