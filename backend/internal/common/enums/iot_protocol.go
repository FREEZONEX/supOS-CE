package enums

import "strings"

// IOTProtocol represents IoT protocols
type IOTProtocol struct {
	Name        string
	DisplayName string
	// ProtocolClass any
}

var (
	IOTProtocolREST     IOTProtocol = IOTProtocol{Name: "rest", DisplayName: "Api"}
	IOTProtocolRelation IOTProtocol = IOTProtocol{Name: "relation", DisplayName: "Relation"}
	IOTProtocolModbus   IOTProtocol = IOTProtocol{Name: "modbus", DisplayName: "Modbus"}
	IOTProtocolMQTT     IOTProtocol = IOTProtocol{Name: "mqtt", DisplayName: "MQTT"}
	IOTProtocolICMP     IOTProtocol = IOTProtocol{Name: "icmp", DisplayName: "ICMP"}
	IOTProtocolOPCUA    IOTProtocol = IOTProtocol{Name: "opcua", DisplayName: "OPC UA"}
	IOTProtocolOPCDA    IOTProtocol = IOTProtocol{Name: "opcda", DisplayName: "OPC DA"}
	IOTProtocolUnknown  IOTProtocol = IOTProtocol{Name: "unknown", DisplayName: "Unknown"}
)

var allIOTProtocols = map[string]IOTProtocol{
	IOTProtocolREST.Name:     IOTProtocolREST,
	IOTProtocolRelation.Name: IOTProtocolRelation,
	IOTProtocolModbus.Name:   IOTProtocolModbus,
	IOTProtocolMQTT.Name:     IOTProtocolMQTT,
	IOTProtocolICMP.Name:     IOTProtocolICMP,
	IOTProtocolOPCUA.Name:    IOTProtocolOPCUA,
	IOTProtocolOPCDA.Name:    IOTProtocolOPCDA,
	IOTProtocolUnknown.Name:  IOTProtocolUnknown,
}

// String returns the string representation
func (p IOTProtocol) String() string {
	return p.Name
}

// GetIOTProtocolByName converts name to IOTProtocol
func GetIOTProtocolByName(name string) IOTProtocol {
	lowerName := strings.ToLower(name)
	protocol, ok := allIOTProtocols[lowerName]
	if ok {
		return protocol
	}
	return IOTProtocolUnknown
}

// Contains checks if the protocol name exists
func IOTProtocolContains(name string) bool {
	return GetIOTProtocolByName(name) != IOTProtocolUnknown
}

// ListSerialProtocol returns a list of serial protocols
func ListSerialProtocol() []map[string]string {
	return []map[string]string{
		{"key": IOTProtocolModbus.Name, "value": IOTProtocolModbus.DisplayName},
		{"key": IOTProtocolOPCUA.Name, "value": IOTProtocolOPCUA.DisplayName},
		{"key": IOTProtocolOPCDA.Name, "value": IOTProtocolOPCDA.DisplayName},
		{"key": IOTProtocolMQTT.Name, "value": IOTProtocolMQTT.DisplayName},
		{"key": IOTProtocolICMP.Name, "value": IOTProtocolICMP.DisplayName},
	}
}
