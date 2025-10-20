package dto

import "time"

// Topology log node constants
const (
	NodePushOriginalData = "pushOriginalData"
	NodePushMqtt         = "pushMqtt"
	NodePullMqtt         = "pullMqtt"
	NodeDataPersistence  = "dataPersistence"
)

// Topology log event code constants
const (
	EventCodeSuccess = "0"
	EventCodeError   = "1"
)

// TopologyNodes is a list of all available topology nodes.
var TopologyNodes = []string{
	NodePushOriginalData,
	NodePushMqtt,
	NodePullMqtt,
	NodeDataPersistence,
}

// TopologyLog is used for logging events in the data processing topology.
type TopologyLog struct {
	UnsId        string `json:"unsId"`
	TopologyNode string `json:"topologyNode"`
	EventCode    string `json:"eventCode"`
	EventMessage string `json:"eventMessage"`
	EventTime    int64  `json:"eventTime"`
}

// NewTopologyLog is a helper function to create a new TopologyLog instance.
// The actual logging to a file/service should be handled by a dedicated logger,
// which would then marshal this struct to JSON.
func NewTopologyLog(unsId, topologyNode, eventCode, eventMessage string) *TopologyLog {
	return &TopologyLog{
		UnsId:        unsId,
		TopologyNode: topologyNode,
		EventCode:    eventCode,
		EventMessage: eventMessage,
		EventTime:    time.Now().UnixMilli(),
	}
}
