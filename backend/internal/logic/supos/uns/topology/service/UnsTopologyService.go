package service

import (
	"backend/internal/common/event"
	"backend/internal/common/event/mount"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/topologylog"
	"backend/internal/config"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

func init() {
	spring.RegisterBean(NewUnsTopologyService())
}

// UnsTopologyService manages topology state and events
type UnsTopologyService struct {
	mu                  sync.RWMutex
	globalTopologyData  *types.GlobalTopologyData
	topologyJson        []byte // JSON representation of globalTopologyData
	globalTopologyDirty bool
	stopChan            chan struct{}
	unsMapper           dao.UnsNamespaceRepo
	httpClient          http.Client
	log                 logx.Logger
	conf                config.Config
	wsService           serviceApi.IWebsocketSender
}

func NewUnsTopologyService() *UnsTopologyService {
	s := &UnsTopologyService{
		globalTopologyDirty: true,
		stopChan:            make(chan struct{}),
		log:                 logx.WithContext(context.Background()),
	}
	s.httpClient.Timeout = 700 * time.Millisecond
	return s
}

const needFlush = runtime.GOOS != "windows"

// startRefreshTask starts background goroutine to refresh topology statistics periodically
func (s *UnsTopologyService) startRefreshTask() {
	if !needFlush {
		logx.Info("Windows environment detected, skipping topology refresh task")
		return
	}

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ticker.Reset(10 * time.Second)
				if s.wsService.HasTopologies() {
					s.refresh()
				} else {
					s.log.Debug("Not refreshing topologies")
				}
			case <-s.stopChan:
				logx.Error("Topology refresh task stopped")
				return
			}
		}
	}()

	logx.Info("Topology refresh task started")
}

// refresh collects topology statistics and updates global data
func (s *UnsTopologyService) refresh() {
	if !needFlush || !s.wsService.HasTopologies() {
		return
	}
	ctx := context.Background()
	db := dao.GetDb(ctx)

	s.mu.Lock()
	defer s.mu.Unlock()

	topologyData := &types.GlobalTopologyData{
		Protocol:    make(map[string]int64),
		ICMPStates:  make([]any, 0),
		MountStatus: make(map[string]string),
	}

	// Count models and instances by path_type
	if pathTypeCounts, err := s.unsMapper.CountByPathType(db); err == nil {
		for _, count := range pathTypeCounts {
			switch count.PathType {
			case 0:
				topologyData.ModelNum = count.Count
			case 2:
				topologyData.InstanceNum = count.Count
			}
			logx.Debugf("PathType=%d, Count=%d", count.PathType, count.Count)
		}
	} else {
		logx.Errorf("failed to count by path type: %v", err)
	}

	// Count protocol instances
	if protocolCounts, err := s.unsMapper.CountByProtocolType(db); err == nil {
		for _, count := range protocolCounts {
			topologyData.Protocol[count.Protocol] = count.Count
			logx.Debugf("Protocol=%s, Count=%d", count.Protocol, count.Count)
		}
	} else {
		logx.Errorf("failed to count by protocol type: %v", err)
	}

	alarmMapper := dao.NewUnsAlarmsDatumRepo(ctx)
	// Count alarms
	if alarmCount, err := alarmMapper.Count(ctx); err == nil {
		topologyData.AlarmNum = alarmCount
		logx.Debugf("Alarm Count=%d", alarmCount)
	} else {
		logx.Errorf("failed to count alarms: %v", err)
	}
	s.statisticsMQTT(topologyData)

	// Query mount status
	topologyData.MountStatus = s.doCountMountStatus(ctx)

	// Preserve ICMP states from previous data
	if s.globalTopologyData != nil && len(s.globalTopologyData.ICMPStates) > 0 {
		topologyData.ICMPStates = s.globalTopologyData.ICMPStates
	}

	s.globalTopologyData = topologyData
	s.globalTopologyDirty = false

	if jsonBytes, err := json.Marshal(topologyData); err == nil {
		s.topologyJson = jsonBytes
	} else {
		logx.Errorf("failed to marshal topology data: %v", err)
		s.topologyJson = []byte("{}")
	}

	spring.PublishEvent(&event.UnsTopologyChangeEvent{TopologyMsg: s.topologyJson})

	logx.Debugf("Topology statistics refreshed: %s", s.topologyJson)
}

type MqttStats struct {
	AllConnections       int64 `json:"connections"`       // Total MQTT connections
	LiveConnections      int64 `json:"live_connections"`  // Active MQTT connections
	MessageInThroughput  int64 `json:"received_msg_rate"` // Messages/sec into MQTT broker
	MessageOutThroughput int64 `json:"sent_msg_rate"`     // Messages/sec out of MQTT broker
}

func (s *UnsTopologyService) statisticsMQTT(data *types.GlobalTopologyData) {
	apiConfig := s.conf.DevLink.Mqtt.OpenApi
	url := apiConfig.Host + "/api/v5/monitor_current"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.log.Error("Error creating request: ", err)
		return
	}
	req.SetBasicAuth(apiConfig.ApiKey, apiConfig.SecretKey)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.log.Error("Error making request: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			s.log.Error("Error reading response body: ", err)
			return
		}
		s.log.Debug("Received mqtt response: ", string(body), url)
		var mqttStats MqttStats
		err = json.Unmarshal(body, &mqttStats)
		if err != nil {
			s.log.Error("Error parsing JSON:", err)
			return
		}
		data.AllConnections = mqttStats.AllConnections
		data.LiveConnections = mqttStats.LiveConnections
		data.MessageInThroughput = mqttStats.MessageInThroughput
		data.MessageOutThroughput = mqttStats.MessageOutThroughput
	} else {
		s.log.Error("HTTP request failed with status code:", resp.StatusCode)
	}
}

// Stop stops the background refresh task
func (s *UnsTopologyService) Stop() {
	close(s.stopChan)
}

// GetLastMsg returns the last topology message as JSON string
func (s *UnsTopologyService) GetLastMsg() []byte {
	if s.topologyJson == nil {
		s.refresh()
	}
	if len(s.topologyJson) > 0 {
		return s.topologyJson
	}
	return []byte("{}")
}

// UpdateTopologyState updates a specific node's state
func (s *UnsTopologyService) UpdateTopologyState(node string, eventCode string) {
	s.globalTopologyDirty = true

	logx.Infof("topology state update requested for node=%s, eventCode=%s", node, eventCode)
	s.refresh()
}

// RefreshTopology triggers an immediate topology refresh
func (s *UnsTopologyService) RefreshTopology() {
	s.globalTopologyDirty = true

	logx.Debug("topology refresh requested")
	s.refresh()
}

// GainTopologyDataOfFile returns topology node states for a specific instance
func (s *UnsTopologyService) GainTopologyDataOfFile(unsId int64) []types.InstanceTopologyData {
	topologyDatas := make([]types.InstanceTopologyData, len(topologylog.TopologyNodes))
	for i, node := range topologylog.TopologyNodes {
		topologyDatas[i] = types.InstanceTopologyData{
			TopologyNode: node,
			EventCode:    topologylog.EventCodeSuccess,
		}
	}
	return topologyDatas
}

// OnEventContextRefreshedEvent handles application startup event
func (s *UnsTopologyService) OnEventContextRefreshedEvent(e *event.ContextRefreshedEvent) error {
	logx.Info("context refreshed, starting topology background task")
	s.conf = e.SvcContext.Config
	s.wsService = spring.GetBean[serviceApi.IWebsocketSender]()
	s.startRefreshTask()
	return nil
}

func (s *UnsTopologyService) OnEventBatchCreateTableEvent(e *event.BatchCreateTableEvent) error {
	logx.Info("namespace change detected, scheduling topology refresh")
	go func() {
		time.Sleep(1 * time.Second)
		s.RefreshTopology()
	}()
	return nil
}
func (s *UnsTopologyService) OnEventRemoveTopicsEvent(e *event.RemoveTopicsEvent) error {
	logx.Info("namespace change detected, scheduling topology refresh")
	go func() {
		time.Sleep(1 * time.Second)
		s.RefreshTopology()
	}()
	return nil
}

// RemoveFromGlobalTopologyData removes a topic from ICMP states
func (s *UnsTopologyService) RemoveFromGlobalTopologyData(topic string) {
	if s.globalTopologyData == nil {
		return
	}

	newStates := make([]any, 0, len(s.globalTopologyData.ICMPStates))
	for _, state := range s.globalTopologyData.ICMPStates {
		newStates = append(newStates, state)
	}
	s.globalTopologyData.ICMPStates = newStates

	if jsonBytes, err := json.Marshal(s.globalTopologyData); err == nil {
		s.topologyJson = jsonBytes
	} else {
		logx.Errorf("failed to marshal topology data: %v", err)
	}
}

// OnEventMountStatusChangeEvent handles mount status change events
func (s *UnsTopologyService) OnEventMountStatusChangeEvent(e *mount.MountStatusChangeEvent) error {
	logx.Info("mount status change detected")

	if s.globalTopologyData != nil {
		ctx := context.Background()
		mountStatus := s.doCountMountStatus(ctx)
		s.globalTopologyData.MountStatus = mountStatus

		if jsonBytes, err := json.Marshal(s.globalTopologyData); err == nil {
			s.topologyJson = jsonBytes
		} else {
			logx.Errorf("failed to marshal topology data: %v", err)
		}
	}

	spring.PublishEvent(&event.UnsTopologyChangeEvent{TopologyMsg: s.topologyJson})
	return nil
}

// doCountMountStatus queries mount status from database
func (s *UnsTopologyService) doCountMountStatus(ctx context.Context) map[string]string {
	mountStatus := make(map[string]string)
	mountMapper := dao.NewUnsMountRepo(ctx)
	mounts, err := mountMapper.FindAll(ctx)
	if err != nil {
		logx.Errorf("failed to query mount status: %v", err)
		return mountStatus
	}
	for _, mount := range mounts {
		mountStatus[mount.TargetAlias] = mount.Status
	}
	return mountStatus
}
