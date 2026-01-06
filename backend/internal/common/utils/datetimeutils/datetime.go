package datetimeutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Go's standard time layouts
	ISO8601Millis = "2006-01-02T15:04:05.000Z07:00"
	CustomFormat1 = "2006-01-02 15:04:05.000-07"
	CustomFormat2 = "2006-01-02 15:04 Z07:00"
	CustomFormat3 = "2006-01-02 15:04:05.00 Z07:00"
	NoTZFormat1   = "2006-01-02 15:04:05.00"
	NoTZFormat2   = "2006-01-02 15:04"
	NoTZFormat3   = "2006-01-02 15"
	ISOLocalDate  = "2006-01-02"
	SimpleSec     = "20060102150405"
)

var (
	utcZone = time.UTC
	// A list of time layouts to try for parsing, ordered by preference
	timeFormats = []string{
		time.RFC3339,
		time.RFC3339Nano,
		ISO8601Millis,
		CustomFormat1,
		CustomFormat2,
		CustomFormat3,
		NoTZFormat1,
		NoTZFormat2,
		NoTZFormat3,
		ISOLocalDate,
		"2006-01-02 15:04:05", // Common format
	}
)

// DateTimeUTC formats milliseconds to a UTC string.
func DateTimeUTC(mills *int64) string {
	if mills == nil || *mills < 1000 {
		return ""
	}
	t := time.UnixMilli(*mills)
	return t.In(utcZone).Format(time.RFC3339)
}

// DateSimple returns the current time in yyyyMMddHHmmss format in local time.
func DateSimple() string {
	return time.Now().Format(SimpleSec)
}

// GetDateTimeStr formats a timestamp (in millis or micros) to a UTC string.
func GetDateTimeStr(inTime any) string {
	if inTime == nil {
		return time.Now().In(utcZone).Format(time.RFC3339)
	}
	ts, ok := inTime.(int64)
	if !ok {
		return ""
	}

	var t time.Time
	if ts > 100000000000000000 { // microseconds
		t = time.Unix(0, ts*1000)
	} else { // milliseconds
		t = time.UnixMilli(ts)
	}
	return t.In(utcZone).Format(time.RFC3339)
}

// ConvertToMillis ensures a timestamp is in milliseconds.
// It truncates longer timestamps and pads shorter ones.
func ConvertToMillis(timestamp int64) int64 {
	s := strconv.FormatInt(timestamp, 10)
	if len(s) > 13 {
		millis, _ := strconv.ParseInt(s[:13], 10, 64)
		return millis
	}
	for len(s) < 13 {
		s += "0"
	}
	millis, _ := strconv.ParseInt(s, 10, 64)
	return millis
}

// ParseDate attempts to parse a datetime string using multiple formats.
// If the string has no timezone, it assumes the local timezone.
func ParseDate(datetime string) (time.Time, error) {
	datetime = strings.TrimSpace(datetime)
	if datetime == "" {
		return time.Time{}, fmt.Errorf("input datetime string is empty")
	}

	for _, layout := range timeFormats {
		// Try parsing with timezone info
		t, err := time.Parse(layout, datetime)
		if err == nil {
			return t, nil
		}
		// Try parsing without timezone info, in local location
		t, err = time.ParseInLocation(layout, datetime, time.Local)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", datetime)
}

// IsValidTime checks if the given value is a valid timestamp or datetime string.
func IsValidTime(value any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return false
		}
		// Check if it's a number (timestamp)
		if _, err := strconv.ParseInt(s, 10, 64); err == nil {
			return true
		}
		// Check against multiple date formats
		if _, err := ParseDate(s); err == nil {
			return true
		}
	}
	return false
}

// ToUTCISO converts a value (timestamp or string) to a UTC ISO string (yyyy-MM-dd'T'HH:mm:ss.SSS'Z').
func ToUTCISO(dateTime any) string {
	if dateTime == nil {
		return ""
	}

	var t time.Time
	var err error

	switch v := dateTime.(type) {
	case int64:
		t = time.UnixMilli(v)
	case int:
		t = time.UnixMilli(int64(v))
	case float64:
		t = time.UnixMilli(int64(v))
	case string:
		s := strings.TrimSpace(v)
		// Try parsing as timestamp first
		if ts, e := strconv.ParseInt(s, 10, 64); e == nil {
			t = time.UnixMilli(ts)
		} else {
			// Try parsing as a date string
			t, err = ParseDate(s)
			if err != nil {
				return s // Return original if parsing fails
			}
		}
	default:
		return fmt.Sprintf("%v", dateTime)
	}

	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}
