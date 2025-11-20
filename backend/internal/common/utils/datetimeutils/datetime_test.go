package datetimeutils

import (
	"fmt"
	"testing"
	"time"
)

func TestDateFmt(t *testing.T) {
	tm := time.Now()
	tm.In(utcZone).Format(time.RFC3339)

	// 示例：2025-11-20 03:33:25.858+00
	targetTime := time.Date(2025, 11, 20, 3, 33, 25, 858000000, time.UTC)

	fmt.Println("直接格式化:", formatTimeDirect(targetTime))
	fmt.Println("UTC格式化:", formatTimeUTC(targetTime))
	fmt.Println("自定义时区格式化:", formatTimeWithCustomTimezone(targetTime))
}

// 方法1：直接使用Format方法
func formatTimeDirect(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000")
}

// 方法2：处理UTC时间的特定格式
func formatTimeUTC(t time.Time) string {
	// 转换为UTC时间
	utcTime := t.UTC()
	return utcTime.Format("2006-01-02 15:04:05.000-07:00")
}

// 方法3：自定义格式化函数，确保+00时区表示
func formatTimeWithCustomTimezone(t time.Time) string {
	// 使用RFC3339变体格式
	formatted := t.Format("2006-01-02 15:04:05.000Z07:00")
	// 将Z替换为+00
	if len(formatted) > 0 && formatted[len(formatted)-6:] == "Z00:00" {
		formatted = formatted[:len(formatted)-6] + "+00"
	}
	return formatted
}
