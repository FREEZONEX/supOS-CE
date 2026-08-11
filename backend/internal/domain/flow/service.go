package flow

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/internal/infra/outbox"
	"backend/internal/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type Service struct {
	flows         *repo.FlowRepo
	outbox        *outbox.Service
	sourceFlowURL string
	eventFlowURL  string
	http          *http.Client
}

type SaveCommand struct {
	ID              int64
	ParentID        int64
	FlowType        string
	NodeType        string
	Name            string
	Description     string
	Template        string
	UnsNodeIDs      []int64
	UserID          int64
	MockData        bool
	MockTopic       string
	MockFields      []MockField
	MockTriggerMode string
	// MockIntervalSeconds 自定义 mock 发送间隔（秒），0 表示默认（10s）
	MockIntervalSeconds int
	// MockIntervalMillis 亚秒级间隔（毫秒），优先于 MockIntervalSeconds
	MockIntervalMillis int
	// MockFunctionJS 自定义 function 节点代码（负责构造 msg.payload 并 return msg；msg.topic 由框架注入）
	MockFunctionJS string
	AutoDeploy     bool
}

type MockField struct {
	Name string
	Type string
}



func New(ctx context.Context, sourceFlowURL, eventFlowURL string) *Service {
	return &Service{
		flows:         repo.NewFlowRepo(ctx),
		outbox:        outbox.New(),
		sourceFlowURL: strings.TrimRight(strings.TrimSpace(sourceFlowURL), "/"),
		eventFlowURL:  strings.TrimRight(strings.TrimSpace(eventFlowURL), "/"),
		http:          &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *Service) List(ctx context.Context, flowType string, parentID int64, keyword string) (map[string]any, error) {
	filterType, err := flowTypeValue(flowType, true)
	if err != nil {
		return nil, err
	}
	flows, err := s.flows.ListFlows(ctx, repo.FlowFilter{FlowType: filterType, ParentID: parentID, Keyword: keyword})
	if err != nil {
		return nil, err
	}
	return s.flowListResp(ctx, flows), nil
}

func (s *Service) Detail(ctx context.Context, id int64) (map[string]any, error) {
	item, err := s.flows.GetFlow(ctx, id)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	return s.flowResp(ctx, item), nil
}

func (s *Service) Create(ctx context.Context, cmd SaveCommand) (map[string]any, error) {
	item, err := commandToFlow(cmd)
	if err != nil {
		return nil, err
	}
	item.CreatedBy = cmd.UserID
	name, err := s.flows.AvailableFlowName(ctx, item.FlowType, item.ParentID, item.Name)
	if err != nil {
		return nil, err
	}
	item.Name = name
	if item.NodeType == 2 {
		item.FlowData = defaultFlowData(item.Name, mockSpecFromCommand(cmd))
	}
	created, err := s.flows.CreateFlow(ctx, item)
	if err != nil {
		return nil, err
	}
	if cmd.AutoDeploy && item.NodeType == 2 {
		return s.Deploy(ctx, created.ID, cmd.UserID, "")
	}
	return s.flowResp(ctx, created), nil
}

func (s *Service) Update(ctx context.Context, cmd SaveCommand) (map[string]any, error) {
	item, err := commandToFlow(cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	item.UpdatedBy = cmd.UserID
	updated, err := s.flows.UpdateFlow(ctx, item)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	return s.flowResp(ctx, updated), nil
}

func (s *Service) Mark(ctx context.Context, id int64, pinned bool, userID int64) (map[string]any, error) {
	item, err := s.flows.MarkFlow(ctx, id, pinned, userID)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	return s.flowResp(ctx, item), nil
}

func (s *Service) BindUns(ctx context.Context, flowID, unsID, userID int64) (map[string]any, error) {
	if flowID <= 0 || unsID <= 0 {
		return nil, ErrInvalid
	}
	current, err := s.flows.GetFlow(ctx, flowID)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	if current.NodeType == 1 || current.FlowType != 1 {
		return nil, ErrInvalid
	}
	for _, nodeID := range current.UnsNodeIDs {
		if nodeID == unsID {
			return s.flowResp(ctx, current), nil
		}
	}
	current.UnsNodeIDs = append(current.UnsNodeIDs, unsID)
	current.UpdatedBy = userID
	updated, err := s.flows.UpdateFlow(ctx, current)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	return s.flowResp(ctx, updated), nil
}

func (s *Service) Delete(ctx context.Context, id, userID int64) (map[string]any, error) {
	current, err := s.flows.GetFlow(ctx, id)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	if err := s.flows.DeleteFlow(ctx, id, userID); err != nil {
		return nil, normalizeNotFound(err)
	}
	if err := s.deleteNodeRedRuntimeFlow(ctx, current.FlowType, current.RuntimeFlowID); err != nil {
		logx.WithContext(ctx).Errorf("delete node-red runtime flow failed flowID=%d runtimeFlowID=%s err=%v", current.ID, current.RuntimeFlowID, err)
	}
	_ = s.outbox.Enqueue(ctx, "flow.deleted.cleanup", "flow", strconv.FormatInt(id, 10), map[string]any{"flowId": id})
	resp := s.flowResp(ctx, current)
	resp["deleted"] = true
	return resp, nil
}

func (s *Service) DeleteByUnsNodeIDs(ctx context.Context, nodeIDs []int64, userID int64) ([]map[string]any, error) {
	if s == nil || len(nodeIDs) == 0 {
		return []map[string]any{}, nil
	}
	ids, err := s.flows.FlowIDsByNodeIDs(ctx, nodeIDs, 1)
	if err != nil {
		return nil, err
	}
	return s.deleteFlowIDs(ctx, ids, userID)
}

func (s *Service) DeleteByUnsNodes(ctx context.Context, nodes []repo.UnsNode, userID int64) ([]map[string]any, error) {
	if s == nil || len(nodes) == 0 {
		return []map[string]any{}, nil
	}
	ids, err := s.flows.FlowIDsByNodeIDs(ctx, unsNodeIDs(nodes), 1)
	if err != nil {
		return nil, err
	}
	return s.deleteFlowIDs(ctx, ids, userID)
}

func unsNodeIDs(nodes []repo.UnsNode) []int64 {
	out := make([]int64, 0, len(nodes))
	for _, node := range nodes {
		if node.ID > 0 {
			out = append(out, node.ID)
		}
	}
	return out
}

func (s *Service) deleteFlowIDs(ctx context.Context, ids []int64, userID int64) ([]map[string]any, error) {
	out := make([]map[string]any, 0, len(ids))
	for _, id := range ids {
		deleted, err := s.Delete(ctx, id, userID)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		out = append(out, deleted)
	}
	return out, nil
}

func (s *Service) SaveData(ctx context.Context, id int64, data string, userID int64) (map[string]any, error) {
	current, err := s.flows.GetFlow(ctx, id)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	if current.NodeType == 1 {
		return nil, repo.ErrFlowFolderCannotDeploy
	}
	normalized, _, err := s.normalizeAndValidateRuntimeNodes(ctx, current, data)
	if err != nil {
		return nil, err
	}
	item, err := s.flows.SaveFlowData(ctx, id, normalized, userID)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	return s.flowResp(ctx, item), nil
}

func (s *Service) Deploy(ctx context.Context, id, userID int64, flowData string) (map[string]any, error) {
	current, err := s.flows.GetFlow(ctx, id)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	if current.NodeType == 1 {
		return nil, repo.ErrFlowFolderCannotDeploy
	}
	targetStatus := "deployed"
	runtimeDisabled := false
	if strings.EqualFold(strings.TrimSpace(current.Status), "disabled") {
		targetStatus = "disabled"
		runtimeDisabled = true
	}
	snapshot, err := s.resolveRuntimeSnapshot(ctx, current, flowData)
	if err != nil {
		return nil, err
	}
	normalized, _, err := s.normalizeAndValidateRuntimeNodes(ctx, current, snapshot)
	if err != nil {
		return nil, err
	}
	current.FlowData = normalized
	runtimeID, err := s.deployRuntime(ctx, current, runtimeDisabled)
	if err != nil {
		return nil, err
	}
	item, err := s.flows.MarkFlowDeployed(ctx, id, runtimeID, current.FlowData, targetStatus, userID)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	_ = s.outbox.Enqueue(ctx, "flow.deploy.runtime", "flow", strconv.FormatInt(id, 10), s.flowResp(ctx, item))
	return s.flowResp(ctx, item), nil
}

func (s *Service) UpdateStatus(ctx context.Context, id, userID int64, status string) (map[string]any, error) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "disabled", "disable":
		current, err := s.flows.GetFlow(ctx, id)
		if err != nil {
			return nil, normalizeNotFound(err)
		}
		if current.NodeType == 1 {
			return nil, repo.ErrFlowFolderCannotDeploy
		}
		if err := s.setRuntimeDisabled(ctx, current, true); err != nil {
			return nil, err
		}
		item, err := s.flows.UpdateFlowStatus(ctx, id, "disabled", userID)
		if err != nil {
			return nil, normalizeNotFound(err)
		}
		_ = s.outbox.Enqueue(ctx, "flow.disable.runtime", "flow", strconv.FormatInt(id, 10), s.flowResp(ctx, item))
		return s.flowResp(ctx, item), nil
	case "deployed", "enabled", "enable", "running":
		return s.deployEnabled(ctx, id, userID)
	default:
		return nil, ErrInvalid
	}
}

func (s *Service) deployEnabled(ctx context.Context, id, userID int64) (map[string]any, error) {
	current, err := s.flows.GetFlow(ctx, id)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	if current.NodeType == 1 {
		return nil, repo.ErrFlowFolderCannotDeploy
	}
	snapshot, err := s.resolveRuntimeSnapshot(ctx, current, "")
	if err != nil {
		return nil, err
	}
	normalized, _, err := s.normalizeAndValidateRuntimeNodes(ctx, current, snapshot)
	if err != nil {
		return nil, err
	}
	current.FlowData = normalized
	runtimeID, err := s.deployRuntime(ctx, current, false)
	if err != nil {
		return nil, err
	}
	item, err := s.flows.MarkFlowDeployed(ctx, id, runtimeID, current.FlowData, "deployed", userID)
	if err != nil {
		return nil, normalizeNotFound(err)
	}
	_ = s.outbox.Enqueue(ctx, "flow.enable.runtime", "flow", strconv.FormatInt(id, 10), s.flowResp(ctx, item))
	return s.flowResp(ctx, item), nil
}






func commandToFlow(cmd SaveCommand) (repo.Flow, error) {
	if strings.TrimSpace(cmd.Name) == "" {
		return repo.Flow{}, ErrInvalid
	}
	flowType, err := flowTypeValue(cmd.FlowType, false)
	if err != nil {
		return repo.Flow{}, err
	}
	nodeType, err := nodeTypeValue(cmd.NodeType)
	if err != nil {
		return repo.Flow{}, err
	}
	if nodeType == 1 && cmd.ParentID != 0 {
		return repo.Flow{}, ErrInvalid
	}
	template := strings.TrimSpace(cmd.Template)
	if nodeType == 2 && template == "" {
		template = "node-red"
	}
	return repo.Flow{
		ParentID:    cmd.ParentID,
		FlowType:    flowType,
		NodeType:    nodeType,
		Name:        strings.TrimSpace(cmd.Name),
		Description: cmd.Description,
		Template:    template,
		UnsNodeIDs:  cmd.UnsNodeIDs,
	}, nil
}

func (s *Service) flowResp(ctx context.Context, item repo.Flow) map[string]any {
	return flowResp(item, s.userNameMap(ctx, item.CreatedBy, item.UpdatedBy))
}

func (s *Service) flowListResp(ctx context.Context, items []repo.Flow) map[string]any {
	userIDs := make([]int64, 0, len(items)*2)
	for _, item := range items {
		userIDs = append(userIDs, item.CreatedBy, item.UpdatedBy)
	}
	userNames := s.userNameMap(ctx, userIDs...)
	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		list = append(list, flowResp(item, userNames))
	}
	return map[string]any{"list": list, "total": len(list)}
}


func (s *Service) userNameMap(ctx context.Context, ids ...int64) map[int64]string {
	userNames, err := s.flows.UserNamesByIDs(ctx, ids)
	if err != nil {
		return map[int64]string{}
	}
	return userNames
}

func userDisplayName(userNames map[int64]string, userID int64) string {
	if userID <= 0 {
		return ""
	}
	if name := strings.TrimSpace(userNames[userID]); name != "" {
		return name
	}
	return strconv.FormatInt(userID, 10)
}

func flowResp(item repo.Flow, userNames map[int64]string) map[string]any {
	creator := userDisplayName(userNames, item.CreatedBy)
	return map[string]any{
		"id":               item.ID,
		"flowId":           item.RuntimeFlowID,
		"runtimeFlowId":    item.RuntimeFlowID,
		"name":             item.Name,
		"flowName":         item.Name,
		"description":      item.Description,
		"flowType":         flowTypeName(item.FlowType),
		"nodeType":         nodeTypeName(item.NodeType),
		"parentId":         item.ParentID,
		"flowData":         item.FlowData,
		"status":           item.Status,
		"flowStatus":       item.Status,
		"template":         item.Template,
		"flowTemplate":     item.Template,
		"unsNodeIds":       item.UnsNodeIDs,
		"sortKey":          item.SortKey,
		"isFavorite":       item.IsFavorite,
		"createdBy":        item.CreatedBy,
		"createdByName":    creator,
		"creator":          creator,
		"operatorName":     creator,
		"updatedBy":        item.UpdatedBy,
		"createdTime":      item.CreatedTime,
		"updatedTime":      item.UpdatedTime,
	}
}


type mockSpec struct {
	enabled         bool
	topic           string
	fields          []MockField
	triggerMode     string
	intervalSeconds int
	intervalMillis  int
	customFunc      string
}

const defaultMQTTBrokerNodeName = "emqx"

const (
	mockTriggerModeAuto   = "auto"
	mockTriggerModeManual = "manual"
)

func flowMQTTBrokerNodeName(string) string { return defaultMQTTBrokerNodeName }

func defaultFlowData(name string, spec mockSpec) string {
	if spec.enabled {
		return mockFlowData(name, spec)
	}
	now := strconv.FormatInt(time.Now().UTC().UnixNano(), 36)
	nodes := []map[string]any{
		{
			"id":       "tab-" + now,
			"type":     "tab",
			"label":    strings.TrimSpace(name),
			"disabled": false,
			"info":     "",
		},
		{
			"id":              "broker-" + now,
			"type":            "mqtt-broker",
			"name":            flowMQTTBrokerNodeName(name),
			"broker":          "emqx",
			"port":            "1883",
			"usetls":          false,
			"protocolVersion": "4",
			"keepalive":       "60",
			"cleansession":    true,
			"birthTopic":      "",
			"birthQos":        "0",
			"birthPayload":    "",
			"closeTopic":      "",
			"closeQos":        "0",
			"closePayload":    "",
			"willTopic":       "",
			"willQos":         "0",
			"willPayload":     "",
			"z":               "",
		},
	}
	data, _ := json.Marshal(nodes)
	return string(data)
}

func mockSpecFromCommand(cmd SaveCommand) mockSpec {
	topic := strings.Trim(cmd.MockTopic, "/")
	fields := normalizeMockFields(cmd.MockFields)
	if len(fields) == 0 {
		fields = []MockField{{Name: "value", Type: "DOUBLE"}}
	}
	return mockSpec{
		enabled:         cmd.MockData,
		topic:           topic,
		fields:          fields,
		triggerMode:     normalizeMockTriggerMode(cmd.MockTriggerMode),
		intervalSeconds: cmd.MockIntervalSeconds,
		intervalMillis:  cmd.MockIntervalMillis,
		customFunc:      strings.TrimSpace(cmd.MockFunctionJS),
	}
}

func mockFlowData(name string, spec mockSpec) string {
	now := strconv.FormatInt(time.Now().UTC().UnixNano(), 36)
	groupID := "group-" + now
	commentID := "comment-" + now
	injectID := "inject-" + now
	funcID := "func-" + now
	mqttID := "mqtt-" + now
	brokerID := "broker-" + now
	topic := strings.Trim(spec.topic, "/")
	if topic == "" {
		topic = strings.Trim(name, "/")
	}
	repeat := "10"
	if spec.intervalSeconds > 0 {
		repeat = strconv.Itoa(spec.intervalSeconds)
	}
	if spec.intervalMillis > 0 {
		repeat = strconv.FormatFloat(float64(spec.intervalMillis)/1000, 'f', -1, 64)
	}
	once := true
	onceDelay := 1
	if spec.triggerMode == mockTriggerModeManual {
		repeat = ""
		once = false
		onceDelay = 0
	}
	nodes := []map[string]any{
		{
			"id":   groupID,
			"type": "group",
			"z":    "",
			"name": "",
			"style": map[string]any{
				"stroke":         "#999999",
				"stroke-opacity": "1",
				"fill":           "none",
				"fill-opacity":   "0",
				"label":          false,
				"color":          "#a4a4a4",
			},
			"nodes": []string{commentID, injectID, funcID, mqttID},
			"x":     100,
			"y":     80,
			"w":     930,
			"h":     250,
		},
		{
			"id":    commentID,
			"type":  "comment",
			"z":     "",
			"g":     groupID,
			"name":  "UNS Source Flow Example",
			"info":  mockFlowCommentInfo(),
			"x":     320,
			"y":     120,
			"w":     340,
			"wires": []any{},
		},
		{
			"id":          injectID,
			"type":        "inject",
			"z":           "",
			"g":           groupID,
			"name":        "",
			"props":       []map[string]any{{"p": "payload"}},
			"repeat":      repeat,
			"crontab":     "",
			"once":        once,
			"onceDelay":   onceDelay,
			"topic":       "",
			"payload":     "",
			"payloadType": "date",
			"x":           320,
			"y":           230,
			"wires":       [][]string{{funcID}},
		},
		{
			"id":         funcID,
			"type":       "function",
			"z":          "",
			"g":          groupID,
			"name":       "mock data",
			"func":       mockFunctionOrCustom(topic, spec),
			"outputs":    1,
			"timeout":    0,
			"noerr":      0,
			"initialize": "",
			"finalize":   "",
			"libs":       []any{},
			"x":          520,
			"y":          230,
			"wires":      [][]string{{mqttID}},
		},
		{
			"id":          mqttID,
			"type":        "mqtt out",
			"z":           "",
			"g":           groupID,
			"name":        "",
			"topic":       "",
			"qos":         "",
			"retain":      "",
			"respTopic":   "",
			"contentType": "",
			"userProps":   "",
			"correl":      "",
			"expiry":      "",
			"broker":      brokerID,
			"x":           750,
			"y":           230,
			"wires":       []any{},
		},
		{
			"id":              brokerID,
			"type":            "mqtt-broker",
			"name":            flowMQTTBrokerNodeName(name),
			"broker":          "emqx",
			"port":            "1883",
			"autoConnect":     true,
			"usetls":          false,
			"protocolVersion": "4",
			"keepalive":       "60",
			"cleansession":    true,
			"autoUnsubscribe": true,
			"birthTopic":      "",
			"birthQos":        "0",
			"birthRetain":     "false",
			"birthPayload":    "",
			"birthMsg":        map[string]any{},
			"closeTopic":      "",
			"closeQos":        "0",
			"closeRetain":     "false",
			"closePayload":    "",
			"closeMsg":        map[string]any{},
			"willTopic":       "",
			"willQos":         "0",
			"willRetain":      "false",
			"willPayload":     "",
			"willMsg":         map[string]any{},
			"userProps":       "",
			"sessionExpiry":   "",
			"z":               "",
		},
	}
	data, _ := json.Marshal(nodes)
	return string(data)
}

func normalizeMockTriggerMode(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case mockTriggerModeManual:
		return mockTriggerModeManual
	default:
		return mockTriggerModeAuto
	}
}

func normalizeMockFields(fields []MockField) []MockField {
	out := make([]MockField, 0, len(fields))
	for _, field := range fields {
		name := strings.TrimSpace(field.Name)
		if name == "" || isReservedMockField(name) {
			continue
		}
		out = append(out, MockField{Name: name, Type: strings.TrimSpace(field.Type)})
	}
	return out
}

func isReservedMockField(name string) bool {
	lower := strings.ToLower(strings.TrimSpace(name))
	switch lower {
	case "_id", "_timestamp", "_quality", "_ct", "_qos", "created_time", "timestamp", "quality":
		return true
	default:
		return strings.EqualFold(name, "timeStamp")
	}
}

func mockFlowCommentInfo() string {
	return "UNS Source Flow Example\n\n" +
		"Click the timestamp node to publish a mock data message to MQTT.\n" +
		"Use this flow to verify the FLOW \u2192 MQTT \u2192 UNS pipeline.\n\n" +
		"To connect a real source:\n" +
		"1. Replace mock data with API / database / MQTT input\n" +
		"2. Update the MQTT topic based on your UNS naming convention\n" +
		"3. Keep the payload in a consistent JSON structure"
}

// mockFunctionOrCustom 自定义函数体优先（框架注入 msg.topic），否则默认随机值生成器
func mockFunctionOrCustom(topic string, spec mockSpec) string {
	if spec.customFunc != "" {
		return "msg.topic = " + jsString(topic) + ";\n" + spec.customFunc
	}
	return mockFunction(topic, spec.fields)
}

func mockFunction(topic string, fields []MockField) string {
	lines := []string{
		"function randomString() {",
		"  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';",
		"  let result = '';",
		"  for (let i = 0; i < 20; i++) result += chars.charAt(Math.floor(Math.random() * chars.length));",
		"  return result;",
		"}",
		"function generateRandomNumber() { return Math.floor(Math.random() * 100); }",
		"function generateRandomFloatWithTwoDecimals() { return Number((Math.random() * 100).toFixed(2)); }",
		"function formatCurDate() { return new Date().toISOString(); }",
		"function getBool() { return generateRandomNumber() > 50; }",
		"msg.topic = " + jsString(topic) + ";",
		"msg.payload = {",
	}
	for i, field := range fields {
		name := strings.TrimSpace(field.Name)
		if name == "" {
			continue
		}
		comma := ","
		if i == len(fields)-1 {
			comma = ""
		}
		lines = append(lines, "  "+jsString(name)+": "+mockValueExpression(field.Type)+comma)
	}
	lines = append(lines, "};", "return msg;")
	return strings.Join(lines, "\n")
}

func mockValueExpression(kind string) string {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "int", "integer", "long":
		return "generateRandomNumber()"
	case "float", "double":
		return "generateRandomFloatWithTwoDecimals()"
	case "bool", "boolean":
		return "getBool()"
	case "datetime", "date", "timestamp":
		return "formatCurDate()"
	default:
		return "randomString()"
	}
}

func jsString(value string) string {
	data, _ := json.Marshal(value)
	return string(data)
}










func (s *Service) deployRuntime(ctx context.Context, flow repo.Flow, disabled bool) (string, error) {
	baseURL := s.runtimeBaseURL(flow.FlowType)
	if baseURL == "" {
		return strings.TrimSpace(flow.RuntimeFlowID), nil
	}
	snapshot, err := s.resolveRuntimeSnapshot(ctx, flow, "")
	if err != nil {
		return "", err
	}
	nodes, err := parseRuntimeNodes(snapshot)
	if err != nil {
		return "", err
	}
	nodes = normalizeSubflowNodes(nodes)
	flowNodes, globalNodes := splitGlobalNodes(nodes)
	runtimeID := strings.TrimSpace(flow.RuntimeFlowID)
	if runtimeID != "" {
		exists, err := s.nodeRedFlowExists(ctx, baseURL, runtimeID)
		if err != nil {
			return "", err
		}
		if !exists {
			runtimeID = ""
		}
	}
	if runtimeID == "" {
		created, err := s.postNodeRed(ctx, baseURL, http.MethodPost, "/flow", map[string]any{
			"id":       "",
			"nodes":    []any{},
			"disabled": false,
			"label":    flow.Name,
			"info":     flow.Description,
		})
		if err != nil {
			return "", err
		}
		runtimeID = strings.TrimSpace(asString(created["id"]))
		if runtimeID == "" {
			return "", errors.New("node-red returned empty flow id")
		}
	}
	if err := s.deployGlobalNodes(ctx, baseURL, globalNodes, referencedNodeIDs(flowNodes)); err != nil {
		return "", err
	}
	runtimeNodes, ownedSubflows, hasSubflow := prepareNodeRedRuntimeNodes(flowNodes, runtimeID)
	if hasSubflow {
		currentNodes, err := s.getNodeRedFlowNodes(ctx, baseURL)
		if err != nil {
			return "", err
		}
		rootTab := nodeRedRootTab(runtimeID, flow)
		rootTab["disabled"] = disabled
		mergedNodes := mergeNodeRedRuntimeNodes(currentNodes, rootTab, runtimeNodes, ownedSubflows)
		if err := s.postNodeRedFlows(ctx, baseURL, mergedNodes); err != nil {
			return "", err
		}
		return runtimeID, nil
	}
	if _, err := s.postNodeRed(ctx, baseURL, http.MethodPut, "/flow/"+urlPath(runtimeID), map[string]any{
		"id":       runtimeID,
		"nodes":    runtimeNodes,
		"disabled": disabled,
		"label":    flow.Name,
		"info":     flow.Description,
	}); err != nil {
		return "", err
	}
	return runtimeID, nil
}

func (s *Service) setRuntimeDisabled(ctx context.Context, flow repo.Flow, disabled bool) error {
	baseURL := s.runtimeBaseURL(flow.FlowType)
	if baseURL == "" {
		return nil
	}
	runtimeID := strings.TrimSpace(flow.RuntimeFlowID)
	if runtimeID == "" {
		return nil
	}
	exists, err := s.nodeRedFlowExists(ctx, baseURL, runtimeID)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	snapshot, err := s.resolveRuntimeSnapshot(ctx, flow, "")
	if err != nil {
		return err
	}
	nodes, err := parseRuntimeNodes(snapshot)
	if err != nil {
		return err
	}
	nodes = normalizeSubflowNodes(nodes)
	flowNodes, _ := splitGlobalNodes(nodes)
	runtimeNodes, ownedSubflows, hasSubflow := prepareNodeRedRuntimeNodes(flowNodes, runtimeID)
	if hasSubflow {
		currentNodes, err := s.getNodeRedFlowNodes(ctx, baseURL)
		if err != nil {
			return err
		}
		rootTab := nodeRedRootTab(runtimeID, flow)
		rootTab["disabled"] = disabled
		mergedNodes := mergeNodeRedRuntimeNodes(currentNodes, rootTab, runtimeNodes, ownedSubflows)
		return s.postNodeRedFlows(ctx, baseURL, mergedNodes)
	}
	_, err = s.postNodeRed(ctx, baseURL, http.MethodPut, "/flow/"+urlPath(runtimeID), map[string]any{
		"id":       runtimeID,
		"nodes":    runtimeNodes,
		"disabled": disabled,
		"label":    flow.Name,
		"info":     flow.Description,
	})
	return err
}

func flowTypeValue(value string, allowEmpty bool) (int, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		if allowEmpty {
			return 0, nil
		}
		return 0, ErrInvalid
	case "source", "1":
		return 1, nil
	case "event", "2":
		return 2, nil
	default:
		return 0, ErrInvalid
	}
}

func flowTypeName(value int) string {
	if value == 2 {
		return "event"
	}
	return "source"
}

func nodeTypeValue(value string) (int, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "flow", "2":
		return 2, nil
	case "folder", "1":
		return 1, nil
	default:
		return 0, ErrInvalid
	}
}

func nodeTypeName(value int) string {
	if value == 1 {
		return "folder"
	}
	return "flow"
}

func normalizeNotFound(err error) error {
	if errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

var (
	ErrInvalid  = errors.New("invalid flow")
	ErrNotFound = errors.New("flow not found")
)

func (s *Service) HasUnsBindings(ctx context.Context, nodes []repo.UnsNode) (bool, error) {
	ids := unsNodeIDs(nodes)
	bound, err := s.flows.FlowIDsByNodeIDs(ctx, ids, 1)
	return len(bound) > 0, err
}
