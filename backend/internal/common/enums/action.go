package enums

// ActionEnum represents action types
type ActionEnum string

const (
	ActionAdd    ActionEnum = "add"
	ActionModify ActionEnum = "modify"
	ActionDelete ActionEnum = "delete"
)

// GetCode returns the action code
func (a ActionEnum) GetCode() string {
	return string(a)
}

// String returns the string representation
func (a ActionEnum) String() string {
	return string(a)
}
