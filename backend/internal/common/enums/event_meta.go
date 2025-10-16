package enums

// EventMetaEnum represents event metadata types
type EventMetaEnum string

const (
	EventMetaUserChange     EventMetaEnum = "user"
	EventMetaRoleChange     EventMetaEnum = "role"
	EventMetaUnsFieldChange EventMetaEnum = "field"
)

// GetCode returns the event meta code
func (e EventMetaEnum) GetCode() string {
	return string(e)
}

// String returns the string representation
func (e EventMetaEnum) String() string {
	return string(e)
}
