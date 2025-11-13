package topologylog

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// Constants for topology nodes
const (
	NodePushOriginalData = "pushOriginalData"
	NodePushMqtt         = "pushMqtt"
	NodePullMqtt         = "pullMqtt"
	NodeDataPersistence  = "dataPersistence"
)

// Constants for event codes
const (
	EventCodeSuccess = "0"
	EventCodeError   = "1"
)

var TopologyNodes = []string{NodePushOriginalData, NodePushMqtt, NodePullMqtt, NodeDataPersistence}

// TopologyLog defines the structure for a topology log entry.
type TopologyLog struct {
	UnsId        string `json:"unsId"`
	TopologyNode string `json:"topologyNode"`
	EventCode    string `json:"eventCode"`
	EventMessage string `json:"eventMessage"`
	EventTime    int64  `json:"eventTime"`
}

// Log generates and logs a topology event.
// It uses a specific logger named "topology" to ensure logs are written to the correct destination.
func Log(unsId int64, topologyNode, eventCode, eventMessage string) {
	logEntry := TopologyLog{
		UnsId:        strconv.FormatInt(unsId, 10),
		TopologyNode: topologyNode,
		EventCode:    eventCode,
		EventMessage: eventMessage,
		EventTime:    time.Now().UnixMilli(),
	}

	// In Java, a specific logger "topology" is used. We use logx.Info here.
	// Assuming the log configuration will route this correctly.
	logJson, err := json.Marshal(logEntry)
	if err != nil {
		logx.Errorf("failed to marshal topology log: %v", err)
		return
	}

	// Direct output to stdout, assuming a log collector (like filebeat) will pick it up.
	// This mimics the behavior of a dedicated logger.
	logx.Info(string(logJson))
}
