package validator

import (
	"backend/internal/common/errors"
	"fmt"
)

// AssertTrue checks if a condition is true, otherwise returns a business error.
func AssertTrue(condition bool, code int, params ...any) error {
	if !condition {
		return errors.NewBuzError(code, fmt.Sprintf("%v", params...))
	}
	return nil
}

// AssertFalse checks if a condition is false, otherwise returns a business error.
func AssertFalse(condition bool, code int, params ...any) error {
	if condition {
		return errors.NewBuzError(code, fmt.Sprintf("%v", params...))
	}
	return nil
}

// AssertNil checks if a value is nil, otherwise returns a business error.
func AssertNil(value any, code int, params ...any) error {
	if value != nil {
		// A more robust nil check might be needed for interfaces
		return errors.NewBuzError(code, fmt.Sprintf("%v", params...))
	}
	return nil
}

// AssertNotNil checks if a value is not nil, otherwise returns a business error.
func AssertNotNil(value any, code int, params ...any) error {
	if value == nil {
		return errors.NewBuzError(code, fmt.Sprintf("%v", params...))
	}
	return nil
}

// AssertNotBlank checks if a string is not blank (empty or whitespace),
// otherwise returns a business error.
func AssertNotBlank(s string, code int, params ...any) error {
	if s == "" { // Simplified from StrUtil.isBlank
		return errors.NewBuzError(code, fmt.Sprintf("%v", params...))
	}
	return nil
}
