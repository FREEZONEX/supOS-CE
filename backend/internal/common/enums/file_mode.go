package enums

import "strings"

// FileReadWriteMode represents file read/write modes
type FileReadWriteMode string

const (
	FileModeReadOnly  FileReadWriteMode = "READ_ONLY"
	FileModeReadWrite FileReadWriteMode = "READ_WRITE"
)

// String returns the string representation
func (frw FileReadWriteMode) String() string {
	return string(frw)
}

// GetFileReadWriteModeByNameIgnoreCase converts string to FileReadWriteMode
func GetFileReadWriteModeByNameIgnoreCase(name string) (FileReadWriteMode, bool) {
	if name == "" {
		return "", false
	}

	upperName := strings.ToUpper(name)
	switch FileReadWriteMode(upperName) {
	case FileModeReadOnly, FileModeReadWrite:
		return FileReadWriteMode(upperName), true
	default:
		return "", false
	}
}
