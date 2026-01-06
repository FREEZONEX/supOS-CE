// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"backend/internal/common/constants"
	"backend/internal/logic/supos/flowcommon"
	"backend/internal/logic/supos/sourceflow"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyEventFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Proxy Node-RED event /flows endpoint using cookie scoped id
func NewProxyEventFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyEventFlowsLogic {
	return &ProxyEventFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyEventFlowsLogic) ProxyEventFlows(flowID string) (string, error) {
	primaryID := strings.TrimSpace(flowID)
	if primaryID == "" {
		return l.fetchAllFlowsFromNodeRed(), nil
	}

	id, err := strconv.ParseInt(primaryID, 10, 64)
	if err != nil || id <= 0 {
		return "", errors.Parameter.WithMsg("nodered.invalid.parameter")
	}

	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	flow, err := sourceflow.LoadFlowByType(l.ctx, repo, id, constants.FlowTypeEVENTFLOW)
	if err != nil {
		return "", err
	}

	nodes, err := l.resolveFlowNodes(flow)
	if err != nil {
		return "", err
	}
	nodes = l.ensureLabelNode(nodes, flow)
	nodes = l.appendGlobalNodes(nodes)

	resp := map[string]any{
		"flows": nodes,
		"rev":   l.fetchVersionRev(),
	}
	data, err := json.Marshal(resp)
	if err != nil {
		l.Errorf("marshal proxy event node-red flow response failed: %v", err)
		return "", errors.System.WithMsg("error.sys.systemError").AddDetail(err.Error())
	}
	return string(data), nil
}

func (l *ProxyEventFlowsLogic) resolveFlowNodes(flow *relationDB.NoderedSourceFlow) ([]map[string]any, error) {
	draft := strings.TrimSpace(flow.FlowData)
	if draft != "" {
		var nodes []map[string]any
		if err := json.Unmarshal([]byte(draft), &nodes); err != nil {
			l.Errorf("unmarshal event flow(%d) draft json failed: %v", flow.ID, err)
			return nil, errors.System.WithMsg("error.sys.systemError").AddDetail(err.Error())
		}
		return nodes, nil
	}

	client := l.svcCtx.EventNodeRed
	flowRuntimeID := strings.TrimSpace(flow.FlowID)
	if client == nil || flowRuntimeID == "" {
		return make([]map[string]any, 0), nil
	}

	var out map[string]any
	code, body, errs := client.GetFlowNodesV1(l.ctx, flowRuntimeID, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("fetch event node-red flow(%s) nodes failed: code=%d err=%v body=%s", flowRuntimeID, code, errs, string(body))
		return make([]map[string]any, 0), nil
	}
	rawNodes, ok := out["nodes"].([]any)
	if !ok || len(rawNodes) == 0 {
		return make([]map[string]any, 0), nil
	}

	nodes := make([]map[string]any, 0, len(rawNodes))
	for _, item := range rawNodes {
		if m, ok := item.(map[string]any); ok && m != nil {
			nodes = append(nodes, m)
		}
	}
	if len(nodes) == 0 {
		return make([]map[string]any, 0), nil
	}
	return nodes, nil
}

func (l *ProxyEventFlowsLogic) ensureLabelNode(nodes []map[string]any, flow *relationDB.NoderedSourceFlow) []map[string]any {
	if nodes == nil {
		nodes = make([]map[string]any, 0)
	}
	tabID := strings.TrimSpace(flow.FlowID)
	if tabID == "" {
		tabID = flowcommon.GenerateNodeID()
	}
	hasTab := false
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if _, ok := node["z"]; ok {
			node["z"] = tabID
		}
		if typ, _ := node["type"].(string); strings.TrimSpace(typ) == "tab" {
			hasTab = true
		}
	}
	if !hasTab {
		labelNode := map[string]any{
			"id":       tabID,
			"type":     "tab",
			"label":    flow.FlowName,
			"disabled": false,
			"info":     flow.Description,
		}
		nodes = append(nodes, labelNode)
	}
	return nodes
}

func (l *ProxyEventFlowsLogic) appendGlobalNodes(nodes []map[string]any) []map[string]any {
	client := l.svcCtx.EventNodeRed
	if client == nil {
		return nodes
	}
	var out map[string]any
	code, body, errs := client.DoJSON(l.ctx, "GET", "/flow/global", nil, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("fetch event global nodes failed: code=%d err=%v body=%s", code, errs, string(body))
		return nodes
	}
	configs, ok := out["configs"].([]any)
	if !ok || len(configs) == 0 {
		return nodes
	}
	existing := make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if id, ok := node["id"].(string); ok {
			existing[strings.TrimSpace(id)] = struct{}{}
		}
	}
	for _, item := range configs {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		id, _ := m["id"].(string)
		id = strings.TrimSpace(id)
		if id != "" {
			if _, exists := existing[id]; exists {
				continue
			}
			existing[id] = struct{}{}
		}
		nodes = append(nodes, m)
	}
	return nodes
}

func (l *ProxyEventFlowsLogic) fetchVersionRev() string {
	client := l.svcCtx.EventNodeRed
	if client == nil {
		return ""
	}
	var out map[string]any
	code, body, errs := client.GetVersionRevV2(l.ctx, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("fetch event node-red flows rev failed: code=%d err=%v body=%s", code, errs, string(body))
		return ""
	}
	if rev, ok := out["rev"].(string); ok {
		return rev
	}
	return ""
}

func (l *ProxyEventFlowsLogic) fetchAllFlowsFromNodeRed() string {
	client := l.svcCtx.EventNodeRed
	if client == nil {
		return "[]"
	}
	var out any
	code, body, errs := client.DoJSON(l.ctx, "GET", "/flows", nil, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("fetch event node-red flows (direct) failed: code=%d err=%v body=%s", code, errs, string(body))
		if len(body) > 0 {
			return string(body)
		}
		return "[]"
	}
	if out == nil {
		return string(body)
	}
	data, err := json.Marshal(out)
	if err != nil {
		l.Errorf("marshal direct event node-red response failed: %v", err)
		return string(body)
	}
	return string(data)
}
