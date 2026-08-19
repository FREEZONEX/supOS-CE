package common

import "unicode/utf8"

const (
	NameMaxLen           = 63
	DescMaxLen           = 512
	AliasMaxLen          = 128
	DisplayNameMaxLen    = 128
	APIKeyNameMaxLen     = 100
	ConnectionNameMaxLen = 128
)

func ExceedMaxLength(value string, maxLength int) bool {
	return utf8.RuneCountInString(value) > maxLength
}

func ExceedNameLimit(name string) bool {
	return ExceedMaxLength(name, NameMaxLen)
}

func ExceedDescLimit(description string) bool {
	return ExceedMaxLength(description, DescMaxLen)
}
