package enums

// ServiceEnum represents service types
type ServiceEnum string

const (
	ServiceAuth ServiceEnum = "auth"
	ServiceUNS  ServiceEnum = "uns"
)

// GetCode returns the service code
func (s ServiceEnum) GetCode() string {
	return string(s)
}

// String returns the string representation
func (s ServiceEnum) String() string {
	return string(s)
}
