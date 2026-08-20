package uns

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"backend/internal/repo"
	"backend/internal/secrets"
	"backend/internal/svc"
)

const (
	globalTopologyRefreshInterval = 10 * time.Second
	emqxMonitorTimeout            = 700 * time.Millisecond
)

type globalTopologyData struct {
	Folder               int64             `json:"Folder"`
	File                 int64             `json:"File"`
	Alarm                int64             `json:"Alarm"`
	AllConnections       int64             `json:"allConnections"`
	LiveConnections      int64             `json:"liveConnections"`
	MessageInThroughput  int64             `json:"messageInThroughput"`
	MessageOutThroughput int64             `json:"messageOutThroughput"`
	Protocol             map[string]int64  `json:"protocol"`
	ICMPStates           []any             `json:"icmpStates"`
	MountStatus          map[string]string `json:"mountStatus"`
}

type emqxMonitorCurrent struct {
	Connections     int64 `json:"connections"`
	LiveConnections int64 `json:"live_connections"`
	ReceivedMsgRate int64 `json:"received_msg_rate"`
	SentMsgRate     int64 `json:"sent_msg_rate"`
}

func streamGlobalTopology(svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request, flusher http.Flusher) {
	sendGlobalTopologySnapshot(svcCtx, w, r.Context(), flusher)

	ticker := time.NewTicker(globalTopologyRefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			sendGlobalTopologySnapshot(svcCtx, w, r.Context(), flusher)
		case <-r.Context().Done():
			return
		}
	}
}

func sendGlobalTopologySnapshot(svcCtx *svc.ServiceContext, w http.ResponseWriter, ctx context.Context, flusher http.Flusher) {
	raw, err := json.Marshal(collectGlobalTopology(ctx, svcCtx))
	if err != nil {
		raw = []byte("{}")
	}
	_, _ = w.Write(append(append([]byte("data: "), raw...), '\n', '\n'))
	if flusher != nil {
		flusher.Flush()
	}
}

func collectGlobalTopology(ctx context.Context, svcCtx *svc.ServiceContext) globalTopologyData {
	data := globalTopologyData{
		Protocol:    map[string]int64{},
		ICMPStates:  []any{},
		MountStatus: map[string]string{},
	}
	collectUnsNodeCounts(ctx, &data)
	collectMqttMonitor(ctx, svcCtx, &data)
	return data
}

func collectUnsNodeCounts(ctx context.Context, data *globalTopologyData) {
	db := repo.GetCommonConn(ctx)
	if db == nil {
		return
	}
	var rows []struct {
		Type  int16 `gorm:"column:type"`
		Count int64 `gorm:"column:count"`
	}
	if err := db.Table("uns_namespace_node_info").
		Select("type, count(*) AS count").
		Where("deleted_time = 0").
		Group("type").
		Scan(&rows).Error; err != nil {
		return
	}
	for _, row := range rows {
		if row.Type == 2 {
			data.File += row.Count
		} else {
			data.Folder += row.Count
		}
	}
}

func collectMqttMonitor(ctx context.Context, svcCtx *svc.ServiceContext, data *globalTopologyData) {
	if svcCtx == nil {
		return
	}
	baseURL := strings.TrimRight(strings.TrimSpace(svcCtx.Config.Gateway.EmqxUrl), "/")
	if baseURL == "" {
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/v5/monitor_current", nil)
	if err != nil {
		return
	}
	apiKey, apiSecret := secrets.EMQXAPIKey()
	req.SetBasicAuth(apiKey, apiSecret)

	resp, err := (&http.Client{Timeout: emqxMonitorTimeout}).Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}
	var stats emqxMonitorCurrent
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return
	}
	data.AllConnections = stats.Connections
	data.LiveConnections = stats.LiveConnections
	data.MessageInThroughput = stats.ReceivedMsgRate
	data.MessageOutThroughput = stats.SentMsgRate
}
