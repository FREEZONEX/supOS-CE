// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	domainflow "backend/internal/domain/flow"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiFlowLegacyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiFlowLegacyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiFlowLegacyLogic {
	return &OpenapiFlowLegacyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiFlowLegacyLogic) List(req *types.OpenapiFlowLegacyListReq) (resp *types.Envelope, err error) {
	flowType, err := normalizeLegacyFlowType(req.FlowType, true)
	if err != nil {
		return nil, logicx.Error(err)
	}
	data, err := l.svcCtx.App.Flow.List(l.ctx, flowType, 0, req.Keyword)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(map[string]any{"list": openapiFlowInfoList(data["list"])}), nil
}

func (l *OpenapiFlowLegacyLogic) Get(req *types.OpenapiFlowLegacyGetReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.Detail(l.ctx, req.Id)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(openapiFlowInfo(data)), nil
}

func (l *OpenapiFlowLegacyLogic) Create(req *types.OpenapiFlowLegacyCreateReq) (resp *types.Envelope, err error) {
	flowType, err := normalizeLegacyFlowType(req.FlowType, false)
	if err != nil {
		return nil, logicx.Error(err)
	}
	data, err := l.svcCtx.App.Flow.Create(l.ctx, domainflow.SaveCommand{
		FlowType:    flowType,
		NodeType:    "flow",
		Name:        req.FlowName,
		Description: req.Description,
		Template:    req.Template,
		UserID:      logicx.UserID(l.ctx),
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(map[string]any{
		"id":       anyInt64(data["id"]),
		"brokerID": injectedMQTTBrokerID(anyString(data["flowData"])),
	}), nil
}

func (l *OpenapiFlowLegacyLogic) Update(req *types.OpenapiFlowLegacyUpdateReq) (resp *types.Envelope, err error) {
	current, err := l.svcCtx.App.Flow.Detail(l.ctx, req.Id)
	if err != nil {
		return nil, logicx.Error(err)
	}
	name := strings.TrimSpace(req.FlowName)
	if name == "" {
		name = anyString(current["flowName"])
	}
	flowType := anyString(current["flowType"])
	_, err = l.svcCtx.App.Flow.Update(l.ctx, domainflow.SaveCommand{
		ID:          req.Id,
		ParentID:    anyInt64(current["parentId"]),
		FlowType:    flowType,
		NodeType:    anyStringOr(current["nodeType"], "flow"),
		Name:        name,
		Description: stringOr(req.Description, anyString(current["description"])),
		Template:    stringOr(req.Template, anyString(current["template"])),
		UnsNodeIDs:  anyInt64Slice(current["unsNodeIds"]),
		UserID:      logicx.UserID(l.ctx),
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(map[string]any{"success": true}), nil
}

func (l *OpenapiFlowLegacyLogic) Delete(req *types.OpenapiFlowLegacyDeleteReq) (resp *types.Envelope, err error) {
	deleted := 0
	for _, id := range req.Ids {
		if id <= 0 {
			continue
		}
		if _, err := l.svcCtx.App.Flow.Delete(l.ctx, id, logicx.UserID(l.ctx)); err != nil {
			return nil, logicx.Error(err)
		}
		deleted++
	}
	_ = deleted
	return respx.Envelope(map[string]any{"success": true}), nil
}

func (l *OpenapiFlowLegacyLogic) FlowData(req *types.OpenapiFlowLegacyGetReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.Detail(l.ctx, req.Id)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(map[string]any{
		"rev":   "",
		"flows": legacyFlowNodes(anyString(data["flowData"])),
	}), nil
}

func (l *OpenapiFlowLegacyLogic) Deploy(req *types.OpenapiFlowLegacyDeployReq) (resp *types.Envelope, err error) {
	userID := logicx.UserID(l.ctx)
	data, err := l.svcCtx.App.Flow.Deploy(l.ctx, req.Id, userID, req.FlowsJson)
	if err != nil {
		return nil, logicx.Error(err)
	}
	flowID := anyString(data["flowId"])
	if flowID == "" {
		flowID = anyString(data["runtimeFlowId"])
	}
	return respx.Envelope(map[string]any{"flowId": flowID}), nil
}

func (l *OpenapiFlowLegacyLogic) Nodes(req *types.OpenapiFlowLegacyNodesReq) (resp *types.Envelope, err error) {
	flowType, err := normalizeLegacyFlowType(req.FlowType, true)
	if err != nil {
		return nil, logicx.Error(err)
	}
	if flowType == "" {
		flowType = "source"
	}
	runtimeData, err := l.svcCtx.App.Flow.NodeTypes(l.ctx, flowType)
	if err != nil {
		return nil, logicx.Error(err)
	}
	nodes, err := openapiRuntimeNodeSets(anyString(runtimeData["nodesJson"]))
	if err != nil {
		return nil, logicx.Error(err)
	}

	return respx.Envelope(map[string]any{
		"nodes": nodes,
	}), nil
}

type openapiRuntimeNodeSet struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Types   []string `json:"types"`
	Enabled *bool    `json:"enabled"`
	Module  string   `json:"module"`
	Version string   `json:"version"`
}

func openapiRuntimeNodeSets(nodesJSON string) ([]map[string]any, error) {
	nodesJSON = strings.TrimSpace(nodesJSON)
	if nodesJSON == "" {
		nodesJSON = "[]"
	}

	var runtimeNodes []openapiRuntimeNodeSet
	if err := json.Unmarshal([]byte(nodesJSON), &runtimeNodes); err != nil {
		return nil, fmt.Errorf("parse Node-RED node types: %w", err)
	}

	nodes := make([]map[string]any, 0, len(runtimeNodes))
	for _, runtimeNode := range runtimeNodes {
		id := strings.TrimSpace(runtimeNode.ID)
		name := strings.TrimSpace(runtimeNode.Name)
		if name == "" {
			name = id
		}
		module := strings.TrimSpace(runtimeNode.Module)
		if module == "" {
			module = id
		}
		enabled := true
		if runtimeNode.Enabled != nil {
			enabled = *runtimeNode.Enabled
		}
		nodes = append(nodes, map[string]any{
			"id":      id,
			"name":    name,
			"types":   runtimeNode.Types,
			"enabled": enabled,
			"module":  module,
			"version": runtimeNode.Version,
		})
	}
	return nodes, nil
}

func openapiFlowInfoList(value any) []map[string]any {
	switch list := value.(type) {
	case []map[string]any:
		out := make([]map[string]any, 0, len(list))
		for _, item := range list {
			out = append(out, openapiFlowInfo(item))
		}
		return out
	case []any:
		out := make([]map[string]any, 0, len(list))
		for _, item := range list {
			if body, ok := item.(map[string]any); ok {
				out = append(out, openapiFlowInfo(body))
			}
		}
		return out
	default:
		return []map[string]any{}
	}
}

func openapiFlowInfo(item map[string]any) map[string]any {
	return map[string]any{
		"id":                 anyInt64(item["id"]),
		"flowId":             anyString(item["flowId"]),
		"flowName":           anyString(item["flowName"]),
		"flowType":           anyString(item["flowType"]),
		"flowStatus":         anyString(item["flowStatus"]),
		"template":           anyString(item["template"]),
		"description":        anyString(item["description"]),
		"isFavorite":         anyInt64(item["isFavorite"]),
		"createdTime":        anyInt64(item["createdTime"]),
		"updatedTime":        anyInt64(item["updatedTime"]),
	}
}

func normalizeLegacyFlowType(value string, allowEmpty bool) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		if allowEmpty {
			return "", nil
		}
	case "source", "sourceflow", "source_flow", "flowtypesource", "node-red", "nodered", "1":
		return "source", nil
	case "event", "eventflow", "event_flow", "flowtypeevent", "event-flow", "2":
		return "event", nil
	}
	return "", domainflow.ErrInvalid
}

func legacyFlowNodes(raw string) []any {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []any{}
	}
	var list []any
	if strings.HasPrefix(raw, "[") {
		if err := json.Unmarshal([]byte(raw), &list); err == nil {
			return list
		}
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return []any{}
	}
	if flows, ok := payload["flows"].([]any); ok {
		return flows
	}
	if nodes, ok := payload["nodes"].([]any); ok {
		return nodes
	}
	return []any{}
}

func injectedMQTTBrokerID(raw string) string {
	for _, item := range legacyFlowNodes(raw) {
		node, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if !isInjectedMQTTBroker(node) {
			continue
		}
		return strings.TrimSpace(anyString(node["id"]))
	}
	return ""
}

func isInjectedMQTTBroker(node map[string]any) bool {
	if strings.TrimSpace(anyString(node["type"])) != "mqtt-broker" {
		return false
	}
	broker := strings.ToLower(strings.TrimSpace(anyString(node["broker"])))
	if broker == "" {
		broker = strings.ToLower(strings.TrimSpace(anyString(node["host"])))
	}
	return broker == "" || broker == "emqx"
}

func anyString(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func anyStringOr(value any, fallback string) string {
	if out := strings.TrimSpace(anyString(value)); out != "" {
		return out
	}
	return fallback
}

func stringOr(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func anyInt64(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case float64:
		return int64(v)
	case json.Number:
		out, _ := v.Int64()
		return out
	default:
		return 0
	}
}

func anyInt64Slice(value any) []int64 {
	switch v := value.(type) {
	case []int64:
		return v
	case []any:
		out := make([]int64, 0, len(v))
		for _, item := range v {
			if id := anyInt64(item); id > 0 {
				out = append(out, id)
			}
		}
		return out
	default:
		return nil
	}
}
