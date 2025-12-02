package types

import (
	"net/http"
	"strings"
)

func GetI18LangByRequest(r *http.Request) string {
	return GetI18Lang(GetAcceptLanguage(r))
}
func GetI18Lang(language string) string {
	if strings.HasPrefix(language, "zh") {
		return "zh_Hans"
	}
	return "en_US"
}
func GetAcceptLanguage(r *http.Request) string {
	acceptLanguage := r.Header.Get("Accept-Language")
	if acceptLanguage == "" {
		return "en" // 默认语言
	}

	// 解析语言偏好，格式如：zh-CN,zh;q=0.9,en;q=0.8
	languages := strings.Split(acceptLanguage, ",")
	if len(languages) > 0 {
		// 取第一个语言（优先级最高）
		primaryLang := strings.Split(languages[0], ";")[0]
		return strings.TrimSpace(primaryLang)
	}

	return "en"
}
