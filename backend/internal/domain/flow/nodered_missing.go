package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type MissingNodeDeleteCommand struct {
	ID       string
	FlowID   string
	Scope    string
	FlowType string
}

const (
	MissingNodeScopeGlobalConfig = "globalConfig"
	MissingNodeScopeFlowNode     = "flowNode"
	MissingNodeScopeFlowConfig   = "flowConfig"
)

type RuntimeMissingNode struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name,omitempty"`
	Scope     string `json:"scope"`
	FlowID    string `json:"flowId,omitempty"`
	FlowLabel string `json:"flowLabel,omitempty"`
	Users     int    `json:"users,omitempty"`
}

type missingNodeDeleteTarget struct {
	ID     string
	FlowID string
	Scope  string
}

func (s *Service) NodeTypes(ctx context.Context, flowType string) (map[string]any, error) {
	flowTypeID, err := flowTypeValue(flowType, false)
	if err != nil {
		return nil, err
	}
	baseURL := s.runtimeBaseURL(flowTypeID)
	if baseURL == "" {
		return map[string]any{"nodesJson": "[]"}, nil
	}
	body, err := s.doNodeRedBytes(ctx, baseURL, http.MethodGet, "/nodes", nil, map[string]string{"Accept": "application/json"})
	if err != nil {
		return nil, err
	}
	nodesJSON := strings.TrimSpace(string(body))
	if nodesJSON == "" {
		nodesJSON = "[]"
	}
	if !json.Valid([]byte(nodesJSON)) {
		return nil, fmt.Errorf("fetch node-red nodes failed: invalid json response")
	}
	return map[string]any{"nodesJson": nodesJSON}, nil
}

func (s *Service) ListMissingNodes(ctx context.Context, flowType string) (map[string]any, error) {
	flowTypeID, err := flowTypeValue(flowType, false)
	if err != nil {
		return nil, err
	}
	nodes, err := s.listMissingRuntimeNodes(ctx, flowTypeID)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		items = append(items, map[string]any{
			"id":        node.ID,
			"type":      node.Type,
			"name":      node.Name,
			"scope":     node.Scope,
			"flowId":    node.FlowID,
			"flowLabel": node.FlowLabel,
			"users":     node.Users,
		})
	}
	return map[string]any{"nodes": items}, nil
}

func (s *Service) DeleteMissingNode(ctx context.Context, cmd MissingNodeDeleteCommand) (map[string]any, error) {
	flowTypeID, err := flowTypeValue(cmd.FlowType, false)
	if err != nil {
		return nil, err
	}
	deleted, err := s.deleteMissingRuntimeNode(ctx, flowTypeID, missingNodeDeleteTarget{
		ID:     cmd.ID,
		FlowID: cmd.FlowID,
		Scope:  cmd.Scope,
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{"deleted": deleted}, nil
}

func (s *Service) CleanupMissingNodes(ctx context.Context, flowType string) (map[string]any, error) {
	flowType = strings.ToLower(strings.TrimSpace(flowType))
	if flowType == "" {
		flowType = "all"
	}
	typeIDs := []int{}
	labels := []string{}
	switch flowType {
	case "all":
		typeIDs = []int{1, 2}
		labels = []string{"source", "event"}
	case "source", "1":
		typeIDs = []int{1}
		labels = []string{"source"}
	case "event", "2":
		typeIDs = []int{2}
		labels = []string{"event"}
	default:
		return nil, ErrInvalid
	}
	total := 0
	byFlowType := make(map[string]int, len(typeIDs))
	for i, flowTypeID := range typeIDs {
		deleted, err := s.cleanupMissingRuntimeNodes(ctx, flowTypeID)
		if err != nil {
			return nil, err
		}
		total += deleted
		byFlowType[labels[i]] = deleted
	}
	return map[string]any{"deleted": total, "byFlowType": byFlowType}, nil
}

func (s *Service) validateNoMissingNodeTypes(ctx context.Context, flowType int, nodes []map[string]any) error {
	if len(nodes) == 0 {
		return nil
	}
	baseURL := s.runtimeBaseURL(flowType)
	if baseURL == "" {
		return nil
	}
	installedTypes := s.fetchInstalledNodeTypes(ctx, baseURL)
	if len(installedTypes) == 0 {
		return nil
	}
	missingTypes := missingNodeTypes(nodes, installedTypes)
	if len(missingTypes) == 0 {
		return nil
	}
	return fmt.Errorf("missing Node-RED node types: %s", strings.Join(missingTypes, ", "))
}

func missingNodeTypes(nodes []map[string]any, installedTypes map[string]bool) []string {
	if len(nodes) == 0 || len(installedTypes) == 0 {
		return nil
	}
	missing := make(map[string]struct{})
	for _, node := range nodes {
		if !isMissingNodeType(node, installedTypes) {
			continue
		}
		if typ := stringField(node, "type"); typ != "" {
			missing[typ] = struct{}{}
		}
	}
	return sortedStrings(missing)
}

func (s *Service) fetchInstalledNodeTypes(ctx context.Context, baseURL string) map[string]bool {
	installed := make(map[string]bool)
	if strings.TrimSpace(baseURL) == "" {
		return installed
	}
	body, err := s.doNodeRedBytes(ctx, baseURL, http.MethodGet, "/nodes", nil, nil)
	if err != nil {
		return installed
	}
	return parseInstalledNodeTypes(body)
}

func parseInstalledNodeTypes(body []byte) map[string]bool {
	installed := make(map[string]bool)
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return installed
	}
	var out []any
	if strings.HasPrefix(trimmed, "[") && json.Unmarshal(body, &out) == nil {
		for _, item := range out {
			nodeSet, ok := item.(map[string]any)
			if !ok {
				continue
			}
			for _, typ := range toAnySlice(nodeSet["types"]) {
				t := strings.TrimSpace(fmt.Sprint(typ))
				if t != "" {
					installed[t] = true
				}
			}
		}
		return installed
	}
	registerTypeRE := regexp.MustCompile(`RED\.nodes\.registerType\(\s*["']([^"']+)["']`)
	for _, match := range registerTypeRE.FindAllStringSubmatch(trimmed, -1) {
		if len(match) < 2 {
			continue
		}
		t := strings.TrimSpace(match[1])
		if t != "" {
			installed[t] = true
		}
	}
	return installed
}

func filterMissingGlobalNodes(nodes []map[string]any, installedTypes map[string]bool) []map[string]any {
	if len(nodes) == 0 || len(installedTypes) == 0 {
		return nodes
	}
	filtered := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		if isMissingGlobalNode(node, installedTypes) {
			continue
		}
		filtered = append(filtered, node)
	}
	return filtered
}

func isMissingOrphanGlobalNode(node map[string]any, installedTypes map[string]bool, allFlows []map[string]any) bool {
	if !isMissingGlobalNode(node, installedTypes) {
		return false
	}
	id := stringField(node, "id")
	for _, item := range allFlows {
		if nodeReferencesID(item, id) {
			return false
		}
	}
	return true
}

func isMissingGlobalNode(node map[string]any, installedTypes map[string]bool) bool {
	if node == nil {
		return false
	}
	if _, hasZ := node["z"]; hasZ {
		return false
	}
	id := stringField(node, "id")
	if id == "" {
		return false
	}
	typ := stringField(node, "type")
	if typ == "" || isBuiltInFlowType(typ) || installedTypes[typ] {
		return false
	}
	return true
}

func isMissingNodeType(node map[string]any, installedTypes map[string]bool) bool {
	if node == nil {
		return false
	}
	id := stringField(node, "id")
	if id == "" {
		return false
	}
	typ := stringField(node, "type")
	if typ == "" || isBuiltInFlowType(typ) || installedTypes[typ] {
		return false
	}
	return true
}

func isBuiltInFlowType(typ string) bool {
	typ = strings.TrimSpace(typ)
	return typ == "tab" || typ == "group" || typ == "subflow" || strings.HasPrefix(typ, "subflow:")
}

func (s *Service) listMissingRuntimeNodes(ctx context.Context, flowType int) ([]RuntimeMissingNode, error) {
	baseURL := s.runtimeBaseURL(flowType)
	if baseURL == "" {
		return nil, nil
	}
	installedTypes := s.fetchInstalledNodeTypes(ctx, baseURL)
	if len(installedTypes) == 0 {
		return nil, nil
	}
	missing := make([]RuntimeMissingNode, 0)
	for _, node := range s.fetchRuntimeGlobalNodes(ctx, baseURL) {
		if isMissingGlobalNode(node, installedTypes) {
			missing = append(missing, runtimeMissingNode(node, MissingNodeScopeGlobalConfig, "global", "On all flows"))
		}
	}
	for _, tab := range s.fetchFlowTabs(ctx, baseURL) {
		flowID := stringField(tab, "id")
		if flowID == "" {
			continue
		}
		out, err := s.getNodeRed(ctx, baseURL, "/flow/"+urlPath(flowID))
		if err != nil {
			return nil, err
		}
		flowLabel := stringField(out, "label")
		if flowLabel == "" {
			flowLabel = stringField(tab, "label")
		}
		for _, node := range toMapSlice(out["nodes"]) {
			if isMissingNodeType(node, installedTypes) {
				missing = append(missing, runtimeMissingNode(node, MissingNodeScopeFlowNode, flowID, flowLabel))
			}
		}
		for _, node := range toMapSlice(out["configs"]) {
			if isMissingNodeType(node, installedTypes) {
				missing = append(missing, runtimeMissingNode(node, MissingNodeScopeFlowConfig, flowID, flowLabel))
			}
		}
	}
	sort.Slice(missing, func(i, j int) bool {
		if missing[i].FlowLabel != missing[j].FlowLabel {
			return missing[i].FlowLabel < missing[j].FlowLabel
		}
		if missing[i].Scope != missing[j].Scope {
			return missing[i].Scope < missing[j].Scope
		}
		return missing[i].ID < missing[j].ID
	})
	return missing, nil
}

func (s *Service) deleteMissingRuntimeNode(ctx context.Context, flowType int, target missingNodeDeleteTarget) (int, error) {
	baseURL := s.runtimeBaseURL(flowType)
	if baseURL == "" {
		return 0, nil
	}
	target.ID = strings.TrimSpace(target.ID)
	target.FlowID = strings.TrimSpace(target.FlowID)
	target.Scope = strings.TrimSpace(target.Scope)
	if target.ID == "" {
		return 0, fmt.Errorf("node-red missing node id is empty")
	}
	installedTypes := s.fetchInstalledNodeTypes(ctx, baseURL)
	if len(installedTypes) == 0 {
		return 0, nil
	}
	if target.FlowID == "" && target.Scope != MissingNodeScopeGlobalConfig {
		resolved, err := s.resolveMissingNodeTarget(ctx, baseURL, installedTypes, target.ID)
		if err != nil {
			return 0, err
		}
		target = resolved
	}
	if target.Scope == MissingNodeScopeGlobalConfig || target.FlowID == "global" {
		return s.deleteMissingGlobalNode(ctx, baseURL, installedTypes, target.ID)
	}
	if target.FlowID == "" {
		return 0, fmt.Errorf("node-red missing node flow id is empty")
	}
	return s.deleteMissingFlowNode(ctx, baseURL, installedTypes, target)
}

func (s *Service) cleanupMissingRuntimeNodes(ctx context.Context, flowType int) (int, error) {
	nodes, err := s.listMissingRuntimeNodes(ctx, flowType)
	if err != nil {
		return 0, err
	}
	deleted := 0
	for _, node := range nodes {
		count, err := s.deleteMissingRuntimeNode(ctx, flowType, missingNodeDeleteTarget{
			ID:     node.ID,
			FlowID: node.FlowID,
			Scope:  node.Scope,
		})
		if err != nil {
			return deleted, err
		}
		deleted += count
	}
	return deleted, nil
}

func (s *Service) deleteMissingGlobalNode(ctx context.Context, baseURL string, installedTypes map[string]bool, nodeID string) (int, error) {
	globalFlow, err := s.getNodeRed(ctx, baseURL, "/flow/global")
	if err != nil {
		return 0, err
	}
	existing := toMapSlice(nodeRedGlobalConfigs(globalFlow))
	if len(existing) == 0 {
		return 0, nil
	}
	filtered := make([]map[string]any, 0, len(existing))
	removed := false
	for _, node := range existing {
		if stringField(node, "id") == nodeID && isMissingGlobalNode(node, installedTypes) {
			removed = true
			continue
		}
		filtered = append(filtered, node)
	}
	if !removed {
		return 0, nil
	}
	err = s.updateNodeRedGlobalConfigs(ctx, baseURL, globalFlow, toAnySlice(filtered))
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *Service) deleteMissingFlowNode(ctx context.Context, baseURL string, installedTypes map[string]bool, target missingNodeDeleteTarget) (int, error) {
	out, err := s.getNodeRed(ctx, baseURL, "/flow/"+urlPath(target.FlowID))
	if err != nil {
		return 0, err
	}
	nodes := toMapSlice(out["nodes"])
	configs := toMapSlice(out["configs"])
	filteredNodes, removedNodeIDs := removeMissingNodeByID(nodes, installedTypes, target.ID)
	filteredConfigs, removedConfigIDs := removeMissingNodeByID(configs, installedTypes, target.ID)
	if len(removedNodeIDs) == 0 && len(removedConfigIDs) == 0 {
		return 0, nil
	}
	removedIDs := make(map[string]struct{}, len(removedNodeIDs)+len(removedConfigIDs))
	for id := range removedNodeIDs {
		removedIDs[id] = struct{}{}
	}
	for id := range removedConfigIDs {
		removedIDs[id] = struct{}{}
	}
	removeWireReferences(filteredNodes, removedIDs)
	_, err = s.postNodeRed(ctx, baseURL, http.MethodPut, "/flow/"+urlPath(target.FlowID), map[string]any{
		"id":       target.FlowID,
		"label":    out["label"],
		"disabled": out["disabled"],
		"info":     out["info"],
		"nodes":    toAnySlice(filteredNodes),
		"configs":  toAnySlice(filteredConfigs),
	})
	if err != nil {
		return 0, err
	}
	return len(removedNodeIDs) + len(removedConfigIDs), nil
}

func removeMissingNodeByID(nodes []map[string]any, installedTypes map[string]bool, nodeID string) ([]map[string]any, map[string]struct{}) {
	filtered := make([]map[string]any, 0, len(nodes))
	removedIDs := make(map[string]struct{})
	for _, node := range nodes {
		if stringField(node, "id") == nodeID && isMissingNodeType(node, installedTypes) {
			removedIDs[nodeID] = struct{}{}
			continue
		}
		filtered = append(filtered, node)
	}
	return filtered, removedIDs
}

func (s *Service) resolveMissingNodeTarget(ctx context.Context, baseURL string, installedTypes map[string]bool, nodeID string) (missingNodeDeleteTarget, error) {
	matches := make([]missingNodeDeleteTarget, 0, 1)
	for _, node := range s.fetchRuntimeGlobalNodes(ctx, baseURL) {
		if stringField(node, "id") == nodeID && isMissingGlobalNode(node, installedTypes) {
			matches = append(matches, missingNodeDeleteTarget{ID: nodeID, FlowID: "global", Scope: MissingNodeScopeGlobalConfig})
		}
	}
	for _, tab := range s.fetchFlowTabs(ctx, baseURL) {
		flowID := stringField(tab, "id")
		if flowID == "" {
			continue
		}
		out, err := s.getNodeRed(ctx, baseURL, "/flow/"+urlPath(flowID))
		if err != nil {
			continue
		}
		for _, node := range toMapSlice(out["nodes"]) {
			if stringField(node, "id") == nodeID && isMissingNodeType(node, installedTypes) {
				matches = append(matches, missingNodeDeleteTarget{ID: nodeID, FlowID: flowID, Scope: MissingNodeScopeFlowNode})
			}
		}
		for _, node := range toMapSlice(out["configs"]) {
			if stringField(node, "id") == nodeID && isMissingNodeType(node, installedTypes) {
				matches = append(matches, missingNodeDeleteTarget{ID: nodeID, FlowID: flowID, Scope: MissingNodeScopeFlowConfig})
			}
		}
	}
	if len(matches) == 0 {
		return missingNodeDeleteTarget{}, fmt.Errorf("node-red missing node does not exist")
	}
	if len(matches) > 1 {
		return missingNodeDeleteTarget{}, fmt.Errorf("node-red missing node target is ambiguous")
	}
	return matches[0], nil
}

func runtimeMissingNode(node map[string]any, scope, flowID, flowLabel string) RuntimeMissingNode {
	id := stringField(node, "id")
	return RuntimeMissingNode{
		ID:        id,
		Type:      stringField(node, "type"),
		Name:      missingNodeDisplayName(node, id),
		Scope:     scope,
		FlowID:    flowID,
		FlowLabel: flowLabel,
		Users:     len(toAnySlice(node["_users"])),
	}
}

func missingNodeDisplayName(node map[string]any, fallback string) string {
	for _, key := range []string{"name", "label", "topic"} {
		name := stringField(node, key)
		if name != "" && name != "<nil>" {
			return name
		}
	}
	return fallback
}

func (s *Service) fetchFlowTabs(ctx context.Context, baseURL string) []map[string]any {
	all, err := s.getNodeRedFlowNodes(ctx, baseURL)
	if err != nil {
		return nil
	}
	tabs := make([]map[string]any, 0)
	for _, node := range all {
		if stringField(node, "type") == "tab" {
			tabs = append(tabs, node)
		}
	}
	return tabs
}

func (s *Service) fetchRuntimeGlobalNodes(ctx context.Context, baseURL string) []map[string]any {
	globalFlow, err := s.getNodeRed(ctx, baseURL, "/flow/global")
	if err != nil {
		return nil
	}
	return toMapSlice(nodeRedGlobalConfigs(globalFlow))
}
