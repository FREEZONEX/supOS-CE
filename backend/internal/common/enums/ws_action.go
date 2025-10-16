package enums

// WSActionEnum represents WebSocket action commands
type WSActionEnum int

const (
	WSCmdSubsEvent    WSActionEnum = 5  // Subscribe to event
	WSCmdUnsubsEvent  WSActionEnum = 4  // Unsubscribe from event
	WSCmdPublishEvent WSActionEnum = 6  // Publish event
	WSCmdResponse     WSActionEnum = 2  // Response
	WSCmdExternTopic  WSActionEnum = 20 // Show external topic list
)

// CmdNo returns the command number
func (w WSActionEnum) CmdNo() int {
	return int(w)
}

// GetWSActionByCmdNo returns WSActionEnum from command number
func GetWSActionByCmdNo(cmdNo int) (WSActionEnum, bool) {
	switch WSActionEnum(cmdNo) {
	case WSCmdSubsEvent, WSCmdUnsubsEvent, WSCmdPublishEvent, WSCmdResponse, WSCmdExternTopic:
		return WSActionEnum(cmdNo), true
	default:
		return 0, false
	}
}
