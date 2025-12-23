package I18nUtils

import (
	"context"

	"gitee.com/unitedrhino/share/i18ns"
)

func GetMessage(k string, args ...any) string {
	if len(args) == 0 {
		return i18ns.LocalizeMsg(k)
	}
	return i18ns.LocalizeMsg(k, args)
}

func GetMessageWithCtx(ctx context.Context, k string, args ...any) string {
	if len(args) == 0 {
		return i18ns.LocalizeMsgWithCtx(ctx, k)
	}
	return i18ns.LocalizeMsgWithCtx(ctx, k, args)
}
func GetMessageWithLang(lang string, k string, args ...any) string {
	if len(args) == 0 {
		return i18ns.LocalizeMsgWithLang(lang, k)
	}
	return i18ns.LocalizeMsgWithLang(lang, k, args)
}
