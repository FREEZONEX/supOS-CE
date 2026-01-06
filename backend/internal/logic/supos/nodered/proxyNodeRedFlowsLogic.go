// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package nodered

import (
	"context"
	"encoding/json"
	"strings"

	"backend/internal/logic/supos/flowcommon"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyNodeRedFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Proxy Node-RED /flows endpoint
func NewProxyNodeRedFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyNodeRedFlowsLogic {
	return &ProxyNodeRedFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyNodeRedFlowsLogic) ProxyNodeRedFlows(flowID string) (string, error) {
	// primaryID := strings.TrimSpace(flowID)
	// if primaryID == "" {
	// 	return l.fetchAllFlowsFromNodeRed(), nil
	// }

	// id, err := strconv.ParseInt(primaryID, 10, 64)
	// if err != nil || id <= 0 {
	// 	return "", errors.Parameter.WithMsg("nodered.invalid.parameter"))
	// }

	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	flow, err := repo.FindOneByFilter(l.ctx, relationDB.NoderedSourceFlowFilter{FlowID: flowID})
	if err != nil {
		return "", err
	}
	if flow == nil {
		return "", errors.NotFind.WithMsg("nodered.flow.not.exist")
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
		l.Errorf("marshal proxy node-red flow response failed: %v", err)
		return "", errors.System.WithMsg("error.sys.systemError").AddDetail(err.Error())
	}
	return string(data), nil
}

func (l *ProxyNodeRedFlowsLogic) resolveFlowNodes(flow *relationDB.NoderedSourceFlow) ([]map[string]any, error) {
	draft := strings.TrimSpace(flow.FlowData)
	if draft != "" {
		var nodes []map[string]any
		if err := json.Unmarshal([]byte(draft), &nodes); err != nil {
			l.Errorf("unmarshal flow(%d) draft json failed: %v", flow.ID, err)
			return nil, errors.System.WithMsg("error.sys.systemError").AddDetail(err.Error())
		}
		return nodes, nil
	}

	client := l.svcCtx.SourceNodeRed
	flowRuntimeID := strings.TrimSpace(flow.FlowID)
	if client == nil || flowRuntimeID == "" {
		return make([]map[string]any, 0), nil
	}

	var out map[string]any
	code, body, errs := client.GetFlowNodesV1(l.ctx, flowRuntimeID, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("fetch node-red flow(%s) nodes failed: code=%d err=%v body=%s", flowRuntimeID, code, errs, string(body))
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

func (l *ProxyNodeRedFlowsLogic) ensureLabelNode(nodes []map[string]any, flow *relationDB.NoderedSourceFlow) []map[string]any {
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

func (l *ProxyNodeRedFlowsLogic) appendGlobalNodes(nodes []map[string]any) []map[string]any {
	client := l.svcCtx.SourceNodeRed
	if client == nil {
		return nodes
	}
	var out map[string]any
	code, body, errs := client.DoJSON(l.ctx, "GET", "/flow/global", nil, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("fetch global nodes failed: code=%d err=%v body=%s", code, errs, string(body))
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

func (l *ProxyNodeRedFlowsLogic) fetchVersionRev() string {
	client := l.svcCtx.SourceNodeRed
	if client == nil {
		return ""
	}
	var out map[string]any
	code, body, errs := client.GetVersionRevV2(l.ctx, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("fetch node-red flows rev failed: code=%d err=%v body=%s", code, errs, string(body))
		return ""
	}
	if rev, ok := out["rev"].(string); ok {
		return rev
	}
	return ""
}

func (l *ProxyNodeRedFlowsLogic) fetchAllFlowsFromNodeRed() string {
	client := l.svcCtx.SourceNodeRed
	if client == nil {
		return "[]"
	}
	var out any
	code, body, errs := client.DoJSON(l.ctx, "GET", "/flows", nil, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("fetch node-red flows (direct) failed: code=%d err=%v body=%s", code, errs, string(body))
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
		l.Errorf("marshal direct node-red response failed: %v", err)
		return string(body)
	}
	return string(data)
}
