package config

import (
	"net/url"
	"strings"
)

const DefaultLoginPath = "/tier0-login"

// NormalizeLoginPath converts configured login path variants like
// "tier0-login", "/tier0-login", and "//tier0-login" into a browser-safe path.
func NormalizeLoginPath(loginPath string) string {
	loginPath = strings.TrimSpace(loginPath)
	if loginPath == "" {
		return DefaultLoginPath
	}

	if parsed, err := url.Parse(loginPath); err == nil && parsed.IsAbs() {
		return loginPath
	}

	path, suffix := splitPathSuffix(loginPath)
	path = strings.TrimSpace(path)
	path = strings.TrimLeft(path, "/")
	if path == "" {
		return DefaultLoginPath + suffix
	}
	return "/" + path + suffix
}

func splitPathSuffix(value string) (string, string) {
	queryIndex := strings.Index(value, "?")
	fragmentIndex := strings.Index(value, "#")
	suffixIndex := -1
	switch {
	case queryIndex >= 0 && fragmentIndex >= 0:
		if queryIndex < fragmentIndex {
			suffixIndex = queryIndex
		} else {
			suffixIndex = fragmentIndex
		}
	case queryIndex >= 0:
		suffixIndex = queryIndex
	case fragmentIndex >= 0:
		suffixIndex = fragmentIndex
	}
	if suffixIndex < 0 {
		return value, ""
	}
	return value[:suffixIndex], value[suffixIndex:]
}
