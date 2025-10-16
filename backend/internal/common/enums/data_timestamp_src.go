package enums

import "strings"

// DataTimestampSrc represents data timestamp source
type DataTimestampSrc string

const (
	DataTimestampFromGateway DataTimestampSrc = "GATEWAY"
	DataTimestampFromServer  DataTimestampSrc = "SERVER"
)

// String returns the string representation
func (dts DataTimestampSrc) String() string {
	return string(dts)
}

// GetDataTimestampSrcFromString converts string to DataTimestampSrc
func GetDataTimestampSrcFromString(typeStr string) (DataTimestampSrc, bool) {
	upperStr := strings.ToUpper(typeStr)
	switch DataTimestampSrc(upperStr) {
	case DataTimestampFromGateway, DataTimestampFromServer:
		return DataTimestampSrc(upperStr), true
	default:
		return "", false
	}
}
