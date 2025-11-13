package service

import (
	"backend/internal/common/event"
	"backend/internal/common/event/mount"
	"backend/internal/common/utils/topologylog"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"encoding/json"
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
	globalTopologyData  *types.GetInstanceTopologyResp
	globalTopologyDirty bool
	stopChan            chan struct{}
	unsMapper           dao.UnsNamespaceRepo
}

func NewUnsTopologyService() *UnsTopologyService {
	s := &UnsTopologyService{
		globalTopologyDirty: true,
		stopChan:            make(chan struct{}),
	}

	// Start background refresh task (similar to Java's init() method)
	s.startRefreshTask()

	return s
}

// startRefreshTask starts a background goroutine to refresh topology periodically
func (s *UnsTopologyService) startRefreshTask() {
	// Skip topology refresh on Windows (for debugging, as per Java code)
	if runtime.GOOS == "windows" {
		logx.Info("Windows environment detected, skipping topology refresh task")
		return
	}

	go func() {
		// Initial delay of 2 seconds
		time.Sleep(2 * time.Second)

		// Refresh immediately
		s.refresh()

		// Then refresh every 10 seconds
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.refresh()
			case <-s.stopChan:
				logx.Info("Topology refresh task stopped")
				return
			}
		}
	}()

	logx.Info("Topology refresh task started (every 10 seconds)")
}

// refresh collects topology statistics and updates global data
func (s *UnsTopologyService) refresh() {
	ctx := context.Background()
	_ = dao.GetDb(ctx) // Will be used for DB queries in TODO sections

	s.mu.Lock()
	defer s.mu.Unlock()

	// Create new topology data
	topologyData := s.createDefaultTopologyData()

	// TODO: Count models and instances
	// Example: modelCount := s.unsMapper.CountByPathType(db, 0)
	// Example: instanceCount := s.unsMapper.CountByPathType(db, 2)

	// TODO: Count protocol instances
	// Example: protocolCounts := s.unsMapper.CountByProtocolType(db)

	// TODO: Count alarms
	// Example: alarmCount := s.alarmMapper.Count(db)

	// TODO: Query EMQX connection count (optional)
	// Example: HTTP request to http://emqx:18083/api/v5/monitor_current

	// Update global data
	s.globalTopologyData = topologyData
	s.globalTopologyDirty = false

	// Publish topology change event
	spring.PublishEvent(&event.UnsTopologyChangeEvent{})

	logx.Debug("Topology statistics refreshed")
}

// Stop stops the background refresh task
func (s *UnsTopologyService) Stop() {
	close(s.stopChan)
}

// GetLastMsg returns the last topology message as JSON string
func (s *UnsTopologyService) GetLastMsg() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.globalTopologyData == nil {
		s.mu.RUnlock()
		s.mu.Lock()
		s.globalTopologyData = s.createDefaultTopologyData()
		s.mu.Unlock()
		s.mu.RLock()
	}

	jsonBytes, err := json.Marshal(s.globalTopologyData)
	if err != nil {
		logx.Errorf("failed to marshal topology data: %v", err)
		return "{}"
	}
	return string(jsonBytes)
}

// UpdateTopologyState updates a specific node's state
func (s *UnsTopologyService) UpdateTopologyState(node string, eventCode string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.globalTopologyData == nil {
		s.globalTopologyData = s.createDefaultTopologyData()
	}

	// Find and update the node
	for i := range s.globalTopologyData.Data {
		if s.globalTopologyData.Data[i].TopologyNode == node {
			s.globalTopologyData.Data[i].EventCode = eventCode
			s.globalTopologyDirty = true

			// Publish topology change event
			spring.PublishEvent(&event.UnsTopologyChangeEvent{})
			break
		}
	}
}

// RefreshTopology marks topology as dirty and triggers refresh
func (s *UnsTopologyService) RefreshTopology() {
	s.mu.Lock()
	s.globalTopologyDirty = true
	s.mu.Unlock()

	// Publish topology change event
	spring.PublishEvent(&event.UnsTopologyChangeEvent{})
}

func (s *UnsTopologyService) createDefaultTopologyData() *types.GetInstanceTopologyResp {
	topologyDatas := make([]types.InstanceTopologyData, len(topologylog.TopologyNodes))
	for i, node := range topologylog.TopologyNodes {
		topologyDatas[i] = types.InstanceTopologyData{
			TopologyNode: node,
			EventCode:    topologylog.EventCodeSuccess,
		}
	}
	return &types.GetInstanceTopologyResp{Data: topologyDatas}
}

// OnEventNamespaceChangeEvent handles namespace change events
// When model/instance changes, refresh topology statistics after 1 second
func (s *UnsTopologyService) OnEventNamespaceChangeEvent(e *event.NamespaceChangeEvent) error {
	logx.Infof("namespace change detected, scheduling topology refresh")

	// Schedule refresh after 1 second (similar to Java's statisticsExecutor.schedule)
	go func() {
		time.Sleep(1 * time.Second)
		s.RefreshTopology()
	}()

	return nil
}

// OnEventMountStatusChangeEvent handles mount status change events
func (s *UnsTopologyService) OnEventMountStatusChangeEvent(e *mount.MountStatusChangeEvent) error {
	logx.Infof("mount status change detected")

	s.mu.Lock()
	// TODO: Update mount status in topology data
	// For now, just mark as dirty
	s.globalTopologyDirty = true
	s.mu.Unlock()

	// Publish topology change event
	spring.PublishEvent(&event.UnsTopologyChangeEvent{})

	return nil
}
