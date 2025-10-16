package enums

// Condition represents comparison conditions
type Condition string

const (
	ConditionGT Condition = ">"  // Greater than
	ConditionGE Condition = ">=" // Greater than or equal
	ConditionLT Condition = "<"  // Less than
	ConditionLE Condition = "<=" // Less than or equal
	ConditionEQ Condition = "="  // Equal
	ConditionNE Condition = "!=" // Not equal
)

// String returns the string representation
func (c Condition) String() string {
	return string(c)
}

func GetConditionValueOfName(name string) (Condition, bool) {
	switch Condition(name) {
	case ConditionGT, ConditionGE, ConditionLT, ConditionLE, ConditionEQ, ConditionNE:
		return Condition(name), true
	default:
		return "", false
	}
}
