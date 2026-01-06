package integerutil

import (
	"fmt"
	"strconv"
	"unicode"
)

// ParseInt parses a string to an *int32. It returns nil if parsing fails.
// This function defaults to parsing in base 10.
func ParseInt(s string) *int32 {
	return ParseIntWithRadix(s, 10)
}

// ParseIntWithRadix parses a string to an *int32 with a specified radix.
// It returns nil if parsing fails, which is the Go equivalent of Java's exception-free
// Integer.parseInt returning null.
func ParseIntWithRadix(s string, radix int) *int32 {
	// strconv.ParseInt is Go's idiomatic way to handle parsing without exceptions.
	val, err := strconv.ParseInt(s, radix, 32)
	if err != nil {
		return nil
	}
	result := int32(val)
	return &result
}

// GetInt dereferences an *int, returning 0 if the pointer is nil.
func GetInt(num *int) int {
	if num != nil {
		return *num
	}
	return 0
}

// GetIntWithDefault dereferences an *int, returning a default value if the pointer is nil.
func GetIntWithDefault(num *int, defaultValue int) int {
	if num != nil {
		return *num
	}
	return defaultValue
}

// GetInt32 dereferences an *int32, returning 0 if the pointer is nil.
func GetInt32(num *int32) int32 {
	if num != nil {
		return *num
	}
	return 0
}

// GetInt32WithDefault dereferences an *int32, returning a default value if the pointer is nil.
func GetInt32WithDefault(num *int32, defaultValue int32) int32 {
	if num != nil {
		return *num
	}
	return defaultValue
}

var noTailNumber = fmt.Errorf("no trailing numbers found")

func ExtractTailNumbers(s string) (int64, error) {
	numbers := ""
	for i := len(s) - 1; i >= 0; i-- {
		if unicode.IsDigit(rune(s[i])) {
			numbers = string(s[i]) + numbers
		} else {
			break
		}
	}

	if numbers == "" {
		return 0, noTailNumber
	}

	return strconv.ParseInt(numbers, 10, 64)
}
