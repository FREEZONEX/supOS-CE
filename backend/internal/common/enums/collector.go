package enums

import "strings"

// CollectorGatewayType represents collector gateway types
type CollectorGatewayType string

const (
	CollectorGatewayGRPC CollectorGatewayType = "GRPC_GATEWAY"
)

// String returns the string representation
func (cgt CollectorGatewayType) String() string {
	return string(cgt)
}

// GetCollectorGatewayTypeFromString converts string to CollectorGatewayType
func GetCollectorGatewayTypeFromString(typeStr string) (CollectorGatewayType, bool) {
	if strings.EqualFold(string(CollectorGatewayGRPC), typeStr) {
		return CollectorGatewayGRPC, true
	}
	return "", false
}

// CollectorGatewayStatus represents collector gateway connection status
type CollectorGatewayStatus string

const (
	CollectorStatusOnline  CollectorGatewayStatus = "online"  // Connected
	CollectorStatusOffline CollectorGatewayStatus = "offline" // Disconnected
)

// String returns the string representation
func (cgs CollectorGatewayStatus) String() string {
	return string(cgs)
}
