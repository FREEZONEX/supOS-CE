package enums

import (
	"strconv"
	"strings"
)

// TimeUnits represents time units
// 时间单位： b（纳秒）、u（微秒）、a（毫秒）、s（秒）、m（分）、h（小时）、d（天）、w（周）
type TimeUnits struct {
	Code     rune
	Multiple int64
}

var (
	TimeUnitNanoSecond  = TimeUnits{Code: 'b', Multiple: 1}
	TimeUnitMicroSecond = TimeUnits{Code: 'u', Multiple: 1000}
	TimeUnitMillsSecond = TimeUnits{Code: 'a', Multiple: 1000000}
	TimeUnitSecond      = TimeUnits{Code: 's', Multiple: 1000000000}
	TimeUnitMinutes     = TimeUnits{Code: 'm', Multiple: 60000000000}
	TimeUnitHours       = TimeUnits{Code: 'h', Multiple: 60000000000 * 60}
	TimeUnitDay         = TimeUnits{Code: 'd', Multiple: 60000000000 * 60 * 24}
	TimeUnitWeek        = TimeUnits{Code: 'w', Multiple: 60000000000 * 60 * 24 * 7}
)

var timeUnitsMap = map[rune]TimeUnits{
	TimeUnitNanoSecond.Code:  TimeUnitNanoSecond,
	TimeUnitMicroSecond.Code: TimeUnitMicroSecond,
	TimeUnitMillsSecond.Code: TimeUnitMillsSecond,
	TimeUnitSecond.Code:      TimeUnitSecond,
	TimeUnitMinutes.Code:     TimeUnitMinutes,
	TimeUnitHours.Code:       TimeUnitHours,
	TimeUnitDay.Code:         TimeUnitDay,
	TimeUnitWeek.Code:        TimeUnitWeek,
}

// TimeUnitsOf returns TimeUnits from code
func TimeUnitsOf(code rune) (TimeUnits, bool) {
	unit, ok := timeUnitsMap[code]
	return unit, ok
}

// TimeUnitsToNanoSecond converts value to nanoseconds
func (tu TimeUnits) TimeUnitsToNanoSecond(value int64) int64 {
	return value * tu.Multiple
}

// TimeUnitsParseToNanoSecond parses string value to nanoseconds
// Example: "10s" -> 10000000000 nanoseconds
func TimeUnitsParseToNanoSecond(value string) (int64, bool) {
	if value == "" {
		return 0, false
	}

	value = strings.TrimSpace(value)
	if len(value) < 2 {
		return 0, false
	}

	// Extract number part
	numStr := strings.TrimSpace(value[:len(value)-1])
	timeNum, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, false
	}

	// Extract unit code
	unitCode := rune(value[len(value)-1])
	unit, ok := TimeUnitsOf(unitCode)
	if !ok {
		return 0, false
	}

	return unit.TimeUnitsToNanoSecond(timeNum), true
}
