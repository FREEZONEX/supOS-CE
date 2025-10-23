package I18nUtils

import "gitee.com/unitedrhino/share/i18ns"

func GetMessage(k string, args ...any) string {
	return i18ns.LocalizeMsg(k, args)
}
