package enums

// IOTDataType represents IoT data types
type IOTDataType string

const (
	IOTDataTypeBasic IOTDataType = "basic"
	IOTDataTypeArray IOTDataType = "array"
	IOTDataTypeJSON  IOTDataType = "key/value(json)"
)

// Description returns the description of the data type
func (t IOTDataType) Description() string {
	return string(t)
}

// String returns the string representation
func (t IOTDataType) String() string {
	return string(t)
}
