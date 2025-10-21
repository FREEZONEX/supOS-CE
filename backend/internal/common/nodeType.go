package common

// NodeType represents node type enumeration
type NodeType int

const (
	NodeTypePath NodeType = 0 + iota
	NodeTypeModel
	NodeTypeInstance
	NodeTypeInstanceForCalc
	NodeTypeInstanceForTimeseries
	NodeTypeAlarmRule
)

// GetCode returns the code for this node type
func (n NodeType) GetCode() int {
	return int(n)
}

// GetNodeTypeValueOf returns NodeType by code
func GetNodeTypeValueOfCode(code int) NodeType {
	switch code {
	case 0:
		return NodeTypePath
	case 1:
		return NodeTypeModel
	case 2:
		return NodeTypeInstance
	case 3:
		return NodeTypeInstanceForCalc
	case 4:
		return NodeTypeInstanceForTimeseries
	case 5:
		return NodeTypeAlarmRule
	default:
		return NodeTypePath
	}
}

// String returns the string representation
func (n NodeType) String() string {
	switch n {
	case NodeTypePath:
		return "Path"
	case NodeTypeModel:
		return "Model"
	case NodeTypeInstance:
		return "Instance"
	case NodeTypeInstanceForCalc:
		return "InstanceForCalc"
	case NodeTypeInstanceForTimeseries:
		return "InstanceForTimeseries"
	case NodeTypeAlarmRule:
		return "AlarmRule"
	default:
		return "Path"
	}
}
