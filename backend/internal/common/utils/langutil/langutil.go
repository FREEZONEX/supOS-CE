package langutil

import (
	"os"
	"strings"
)

const defaultSystemLocale = "zh-CN"

// SystemLocale returns the current system locale based on SYS_OS_LANG env or the default.
func SystemLocale() string {
	if val := strings.TrimSpace(os.Getenv("SYS_OS_LANG")); val != "" {
		return val
	}
	return defaultSystemLocale
}
