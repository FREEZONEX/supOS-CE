package errutil

import (
	"fmt"
)

// SetErrMsg creates a formatted validation error message string.
// It matches the name and parameters of the original Java `ConstraintErrTipUtils.setErrMsg` method.
// The `ConstraintValidatorContext` parameter from the Java version is omitted
// as there is no direct equivalent in Go's standard validation approach. Instead, this function
// returns the formatted string to be used by the caller, typically to create an error object.
//
// TODO: Integrate with the I18n service to translate the messageKey.
func SetErrMsg(propertyPath, messageKey string, args ...any) string {
	// In the future, this part will call the i18n service.
	// translatedMessage := i18n.GetMessage(messageKey, args...)
	// For now, we just format it.
	var translatedMessage string
	if len(args) > 0 {
		translatedMessage = fmt.Sprintf(messageKey, args...)
	} else {
		translatedMessage = messageKey
	}

	return fmt.Sprintf("%s: %s", propertyPath, translatedMessage)
}
