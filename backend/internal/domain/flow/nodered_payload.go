package flow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"backend/internal/repo"
	"backend/internal/secrets"
)

const nodeRedInternalTokenHeader = "X-Tier0-Internal-Token"

type EditorPayload struct {
	Flows       []map[string]any             `json:"flows"`
	Rev         string                       `json:"rev"`
	Credentials map[string]map[string]string `json:"credentials"`
}

func EmptyEditorPayload() EditorPayload {
	return EditorPayload{
		Flows:       []map[string]any{},
		Rev:         "",
		Credentials: map[string]map[string]string{},
	}
}

func (s *Service) EditorPayload(ctx context.Context, flowType, id int64) (EditorPayload, error) {
	if id <= 0 {
		return EmptyEditorPayload(), ErrInvalid
	}
	current, err := s.flows.GetFlow(ctx, id)
	if err != nil {
		return EmptyEditorPayload(), normalizeNotFound(err)
	}
	if current.FlowType != int(flowType) || current.NodeType == 1 {
		return EmptyEditorPayload(), ErrNotFound
	}
	nodes, err := s.resolveEditorNodes(ctx, current)
	if err != nil {
		return EmptyEditorPayload(), err
	}
	nodes = orderNodesForSubflowLoad(nodes)
	nodes = ensureNodeRedEditorTab(nodes, current.RuntimeFlowID, current.Name, current.Description)
	applyNodeRedEditorDisabled(nodes, strings.EqualFold(strings.TrimSpace(current.Status), "disabled"))
	nodes = s.addEditorLinkContextNodes(ctx, current, nodes)
	nodes = s.addRuntimeGlobalNodes(ctx, current.FlowType, nodes)
	return EditorPayload{
		Flows:       nodes,
		Rev:         s.editorRevision(ctx, current),
		Credentials: nodeRedCredentialsFromNodes(nodes),
	}, nil
}

func (s *Service) resolveEditorNodes(ctx context.Context, flow repo.Flow) ([]map[string]any, error) {
	if strings.TrimSpace(flow.FlowData) != "" {
		nodes, err := parseRuntimeNodes(flow.FlowData)
		if err != nil {
			return nil, err
		}
		if !hasMissingSubflowDefinitions(nodes) {
			return orderNodesForSubflowLoad(nodes), nil
		}
		return normalizeSubflowNodes(nodes), nil
	}
	return s.fetchRuntimeFlowNodes(ctx, flow.FlowType, flow.RuntimeFlowID)
}

func (s *Service) resolveRuntimeSnapshot(ctx context.Context, flow repo.Flow, override string) (string, error) {
	if raw := strings.TrimSpace(override); raw != "" {
		return raw, nil
	}
	if raw := strings.TrimSpace(flow.FlowData); raw != "" {
		return raw, nil
	}
	runtimeNodes, err := s.fetchRuntimeFlowNodes(ctx, flow.FlowType, flow.RuntimeFlowID)
	if err != nil {
		return "", err
	}
	return marshalRuntimeNodes(runtimeNodes)
}


func (s *Service) addEditorLinkContextNodes(ctx context.Context, current repo.Flow, nodes []map[string]any) []map[string]any {
	flows, err := s.flows.ListEditorContextFlows(ctx, current.FlowType, current.ID)
	if err != nil || len(flows) == 0 {
		return nodes
	}
	existingIDs := nodeRedNodeIDSet(nodes)
	out := append([]map[string]any{}, nodes...)
	for _, flow := range flows {
		otherNodes, err := s.resolveEditorNodes(ctx, flow)
		if err != nil {
			continue
		}
		otherNodes = orderNodesForSubflowLoad(otherNodes)
		otherNodes = ensureNodeRedEditorTab(otherNodes, flow.RuntimeFlowID, flow.Name, flow.Description)
		contextNodes := nodeRedLinkContextNodes(flow, otherNodes, existingIDs)
		if len(contextNodes) == 0 {
			continue
		}
		out = append(out, contextNodes...)
	}
	return out
}

func nodeRedNodeIDSet(nodes []map[string]any) map[string]struct{} {
	out := make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		if id := stringField(node, "id"); id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

func nodeRedLinkContextNodes(flow repo.Flow, nodes []map[string]any, existingIDs map[string]struct{}) []map[string]any {
	linkNodes := make([]map[string]any, 0)
	neededContainers := make(map[string]struct{})
	for _, node := range nodes {
		if !isNodeRedLinkNode(node) {
			continue
		}
		z := stringField(node, "z")
		if z == "" {
			continue
		}
		linkNodes = append(linkNodes, node)
		neededContainers[z] = struct{}{}
	}
	if len(linkNodes) == 0 {
		return nil
	}

	out := make([]map[string]any, 0, len(linkNodes)+len(neededContainers))
	availableContainers := make(map[string]struct{}, len(neededContainers))
	for _, node := range nodes {
		typ := stringField(node, "type")
		if typ != "tab" && typ != "subflow" {
			continue
		}
		id := stringField(node, "id")
		if _, needed := neededContainers[id]; !needed {
			continue
		}
		if id == "" {
			continue
		}
		if _, exists := existingIDs[id]; exists {
			continue
		}
		existingIDs[id] = struct{}{}
		availableContainers[id] = struct{}{}
		out = append(out, markNodeRedContextNode(flow, node))
	}
	for _, node := range linkNodes {
		z := stringField(node, "z")
		if _, ok := availableContainers[z]; !ok {
			continue
		}
		id := stringField(node, "id")
		if id == "" {
			continue
		}
		if _, exists := existingIDs[id]; exists {
			continue
		}
		existingIDs[id] = struct{}{}
		out = append(out, markNodeRedContextNode(flow, node))
	}
	return out
}

func isNodeRedLinkNode(node map[string]any) bool {
	switch stringField(node, "type") {
	case "link in", "link out", "link call":
		return true
	default:
		return false
	}
}

func markNodeRedContextNode(flow repo.Flow, node map[string]any) map[string]any {
	out := cloneNodeRedNode(node)
	out["_contextOnly"] = true
	out["_contextFlowId"] = strconv.FormatInt(flow.ID, 10)
	out["_contextFlowName"] = flow.Name
	out["_contextRootWorkspaceId"] = strings.TrimSpace(flow.RuntimeFlowID)
	return out
}

func cloneNodeRedNode(node map[string]any) map[string]any {
	out := make(map[string]any, len(node)+4)
	for key, value := range node {
		out[key] = value
	}
	return out
}

func (s *Service) fetchRuntimeFlowNodes(ctx context.Context, flowType int, runtimeID string) ([]map[string]any, error) {
	runtimeID = strings.TrimSpace(runtimeID)
	if runtimeID == "" {
		return []map[string]any{}, nil
	}
	baseURL := s.runtimeBaseURL(flowType)
	if baseURL == "" {
		return []map[string]any{}, nil
	}
	out, err := s.getNodeRed(ctx, baseURL, "/flow/"+urlPath(runtimeID))
	if err != nil {
		return nil, err
	}
	nodes := append([]map[string]any{}, toMapSlice(out["nodes"])...)
	nodes = append(nodes, toMapSlice(out["configs"])...)
	if len(nodes) == 0 {
		return []map[string]any{}, nil
	}
	nodes, err = s.completeReferencedSubflowsFromRuntime(ctx, baseURL, nodes)
	if err != nil {
		return nil, err
	}
	return orderNodesForSubflowLoad(nodes), nil
}

func (s *Service) completeReferencedSubflowsFromRuntime(ctx context.Context, baseURL string, nodes []map[string]any) ([]map[string]any, error) {
	refs := collectSubflowReferenceIDs(nodes)
	if len(refs) == 0 {
		return nodes, nil
	}
	all, err := s.getNodeRedFlowNodes(ctx, baseURL)
	if err != nil {
		return nil, err
	}
	return appendSubflowClosure(nodes, all, refs), nil
}

func appendSubflowClosure(base []map[string]any, all []map[string]any, refs map[string]struct{}) []map[string]any {
	byID := make(map[string]map[string]any, len(all))
	childrenByZ := make(map[string][]map[string]any)
	existing := make(map[string]struct{}, len(base))
	for _, node := range base {
		if id := stringField(node, "id"); id != "" {
			existing[id] = struct{}{}
		}
	}
	for _, node := range all {
		if id := stringField(node, "id"); id != "" {
			byID[id] = node
		}
		if z := stringField(node, "z"); z != "" {
			childrenByZ[z] = append(childrenByZ[z], node)
		}
	}
	out := append([]map[string]any{}, base...)
	queue := make([]string, 0, len(refs))
	seen := make(map[string]struct{}, len(refs))
	for id := range refs {
		queue = append(queue, id)
	}
	for len(queue) > 0 {
		subflowID := queue[0]
		queue = queue[1:]
		if _, ok := seen[subflowID]; ok {
			continue
		}
		seen[subflowID] = struct{}{}
		if def := byID[subflowID]; stringField(def, "type") == "subflow" {
			if _, ok := existing[subflowID]; !ok {
				out = append(out, def)
				existing[subflowID] = struct{}{}
			}
		}
		for _, child := range childrenByZ[subflowID] {
			if id := stringField(child, "id"); id != "" {
				if _, ok := existing[id]; !ok {
					out = append(out, child)
					existing[id] = struct{}{}
				}
			}
			for nestedID := range collectSubflowReferenceIDs([]map[string]any{child}) {
				if _, ok := seen[nestedID]; !ok {
					queue = append(queue, nestedID)
				}
			}
		}
	}
	return orderNodesForSubflowLoad(out)
}

func (s *Service) editorRevision(ctx context.Context, flow repo.Flow) string {
	if rev, err := s.runtimeRevision(ctx, flow.FlowType); err == nil && strings.TrimSpace(rev) != "" {
		return strings.TrimSpace(rev)
	}
	if strings.TrimSpace(flow.RuntimeFlowID) != "" {
		return strings.TrimSpace(flow.RuntimeFlowID)
	}
	if flow.ID > 0 {
		return strconv.FormatInt(flow.ID, 10)
	}
	return ""
}

func ensureNodeRedEditorTab(nodes []map[string]any, tabID, label, info string) []map[string]any {
	tabID = strings.TrimSpace(tabID)
	if tabID == "" {
		tabID = "flow-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	}
	subflowIDs := nodeRedSubflowIDs(nodes)
	hasTab := false
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if _, ok := node["z"]; ok {
			z := stringField(node, "z")
			if _, isSubflow := subflowIDs[z]; !isSubflow {
				node["z"] = tabID
			}
		}
		if stringField(node, "type") == "tab" {
			hasTab = true
			node["id"] = tabID
			node["label"] = label
			node["info"] = info
			node["disabled"] = false
		}
	}
	if !hasTab {
		nodes = append([]map[string]any{{
			"id":       tabID,
			"type":     "tab",
			"label":    label,
			"disabled": false,
			"info":     info,
		}}, nodes...)
	}
	return nodes
}

func applyNodeRedEditorDisabled(nodes []map[string]any, disabled bool) {
	for _, node := range nodes {
		if node != nil && stringField(node, "type") == "tab" {
			node["disabled"] = disabled
		}
	}
}

func nodeRedSubflowIDs(nodes []map[string]any) map[string]struct{} {
	ids := make(map[string]struct{})
	for _, node := range nodes {
		if node == nil || stringField(node, "type") != "subflow" {
			continue
		}
		if id := stringField(node, "id"); id != "" {
			ids[id] = struct{}{}
		}
	}
	return ids
}

func (s *Service) addRuntimeGlobalNodes(ctx context.Context, flowType int, nodes []map[string]any) []map[string]any {
	baseURL := s.runtimeBaseURL(flowType)
	if baseURL == "" {
		return nodes
	}
	globalNodes := s.fetchRuntimeGlobalNodes(ctx, baseURL)
	if len(globalNodes) == 0 {
		return nodes
	}
	installedTypes := s.fetchInstalledNodeTypes(ctx, baseURL)
	referencedIDs := referencedNodeIDs(nodes)
	exist := make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		if id := stringField(node, "id"); id != "" {
			exist[id] = struct{}{}
		}
	}
	for _, global := range globalNodes {
		id := stringField(global, "id")
		if id == "" || id == "<nil>" {
			nodes = append(nodes, global)
			continue
		}
		if len(installedTypes) > 0 && isMissingGlobalNode(global, installedTypes) {
			if _, referenced := referencedIDs[id]; !referenced {
				continue
			}
		}
		if _, ok := exist[id]; ok {
			continue
		}
		nodes = append(nodes, global)
		exist[id] = struct{}{}
	}
	return nodes
}

func parseRuntimeNodes(raw string) ([]map[string]any, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return []map[string]any{}, nil
	}
	var nodes []map[string]any
	if strings.HasPrefix(raw, "[") {
		if err := json.Unmarshal([]byte(raw), &nodes); err != nil {
			return nil, err
		}
		return nodes, nil
	}
	var payload struct {
		Nodes []map[string]any `json:"nodes"`
		Flows []map[string]any `json:"flows"`
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, err
	}
	if payload.Nodes != nil {
		return payload.Nodes, nil
	}
	if payload.Flows != nil {
		return payload.Flows, nil
	}
	return []map[string]any{}, nil
}

func normalizeRuntimeNodesJSON(raw string) (string, []map[string]any, error) {
	nodes, err := parseRuntimeNodes(raw)
	if err != nil {
		return "", nil, err
	}
	nodes = normalizeSubflowNodes(nodes)
	data, err := marshalRuntimeNodes(nodes)
	if err != nil {
		return "", nil, err
	}
	return data, nodes, nil
}

func marshalRuntimeNodes(nodes []map[string]any) (string, error) {
	if nodes == nil {
		nodes = []map[string]any{}
	}
	data, err := json.Marshal(nodes)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func hasMissingSubflowDefinitions(nodes []map[string]any) bool {
	defs := make(map[string]struct{})
	for _, node := range nodes {
		if node == nil || stringField(node, "type") != "subflow" {
			continue
		}
		if id := stringField(node, "id"); id != "" {
			defs[id] = struct{}{}
		}
	}
	for id := range collectSubflowReferenceIDs(nodes) {
		if _, ok := defs[id]; !ok {
			return true
		}
	}
	return false
}

func normalizeSubflowNodes(nodes []map[string]any) []map[string]any {
	return orderNodesForSubflowLoad(pruneMissingSubflowReferences(nodes))
}

func pruneMissingSubflowReferences(nodes []map[string]any) []map[string]any {
	defs := make(map[string]struct{})
	for _, node := range nodes {
		if node == nil || stringField(node, "type") != "subflow" {
			continue
		}
		if id := stringField(node, "id"); id != "" {
			defs[id] = struct{}{}
		}
	}
	missingSubflows := make(map[string]struct{})
	for id := range collectSubflowReferenceIDs(nodes) {
		if _, ok := defs[id]; !ok {
			missingSubflows[id] = struct{}{}
		}
	}
	if len(missingSubflows) == 0 {
		return nodes
	}
	removedNodeIDs := make(map[string]struct{})
	kept := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		typ := stringField(node, "type")
		if strings.HasPrefix(typ, "subflow:") {
			refID := strings.TrimSpace(strings.TrimPrefix(typ, "subflow:"))
			if _, missing := missingSubflows[refID]; missing {
				if id := stringField(node, "id"); id != "" {
					removedNodeIDs[id] = struct{}{}
				}
				continue
			}
		}
		if z := stringField(node, "z"); z != "" {
			if _, missing := missingSubflows[z]; missing {
				if id := stringField(node, "id"); id != "" {
					removedNodeIDs[id] = struct{}{}
				}
				continue
			}
		}
		kept = append(kept, node)
	}
	removeWireReferences(kept, removedNodeIDs)
	return kept
}

func orderNodesForSubflowLoad(nodes []map[string]any) []map[string]any {
	if len(nodes) == 0 {
		return nodes
	}
	defsByID := make(map[string]map[string]any)
	childrenByZ := make(map[string][]map[string]any)
	for _, node := range nodes {
		id := stringField(node, "id")
		if stringField(node, "type") == "subflow" && id != "" {
			defsByID[id] = node
		}
		if z := stringField(node, "z"); z != "" {
			childrenByZ[z] = append(childrenByZ[z], node)
		}
	}
	out := make([]map[string]any, 0, len(nodes))
	added := make(map[string]struct{}, len(nodes))
	addNode := func(node map[string]any) {
		id := stringField(node, "id")
		if id != "" {
			if _, ok := added[id]; ok {
				return
			}
			added[id] = struct{}{}
		}
		out = append(out, node)
	}
	var addSubflow func(id string)
	visiting := make(map[string]struct{})
	addSubflow = func(id string) {
		if id == "" {
			return
		}
		if _, ok := added[id]; ok {
			return
		}
		if _, ok := visiting[id]; ok {
			return
		}
		def := defsByID[id]
		if def == nil {
			return
		}
		visiting[id] = struct{}{}
		addNode(def)
		for _, child := range childrenByZ[id] {
			for nestedID := range collectSubflowReferenceIDs([]map[string]any{child}) {
				addSubflow(nestedID)
			}
			addNode(child)
		}
		delete(visiting, id)
	}
	for _, node := range nodes {
		if stringField(node, "type") == "tab" {
			addNode(node)
		}
	}
	for _, node := range nodes {
		if stringField(node, "type") == "subflow" {
			addSubflow(stringField(node, "id"))
		}
	}
	for _, node := range nodes {
		for refID := range collectSubflowReferenceIDs([]map[string]any{node}) {
			addSubflow(refID)
		}
		addNode(node)
	}
	return out
}

func collectSubflowReferenceIDs(nodes []map[string]any) map[string]struct{} {
	ids := make(map[string]struct{})
	for _, node := range nodes {
		typ := stringField(node, "type")
		if !strings.HasPrefix(typ, "subflow:") {
			continue
		}
		id := strings.TrimSpace(strings.TrimPrefix(typ, "subflow:"))
		if id != "" {
			ids[id] = struct{}{}
		}
	}
	return ids
}

func splitGlobalNodes(nodes []map[string]any) (flowNodes []map[string]any, globalNodes []map[string]any) {
	flowNodes = make([]map[string]any, 0, len(nodes))
	globalNodes = make([]map[string]any, 0)
	for _, node := range nodes {
		typ := stringField(node, "type")
		_, hasZ := node["z"]
		if typ != "tab" && typ != "subflow" && !hasZ {
			globalNodes = append(globalNodes, node)
			continue
		}
		flowNodes = append(flowNodes, node)
	}
	return flowNodes, globalNodes
}

func (s *Service) normalizeAndValidateRuntimeNodes(ctx context.Context, flow repo.Flow, raw string) (string, []map[string]any, error) {
	nodes, err := parseRuntimeNodes(raw)
	if err != nil {
		return "", nil, err
	}
	nodes = stripNodeRedContextNodes(nodes)
	nodes = filterNodeRedRuntimeScope(nodes, flow.RuntimeFlowID)
	nodes = normalizeSubflowNodes(nodes)
	data, err := marshalRuntimeNodes(nodes)
	if err != nil {
		return "", nil, err
	}
	if err := s.validateNoMissingNodeTypes(ctx, flow.FlowType, nodes); err != nil {
		return "", nil, err
	}
	return data, nodes, nil
}

func stripNodeRedContextNodes(nodes []map[string]any) []map[string]any {
	out := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		if isNodeRedContextNode(node) {
			continue
		}
		out = append(out, node)
	}
	return out
}

func isNodeRedContextNode(node map[string]any) bool {
	value, ok := node["_contextOnly"]
	if !ok {
		return false
	}
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(strings.TrimSpace(v), "true")
	default:
		return fmt.Sprint(v) == "true"
	}
}

func filterNodeRedRuntimeScope(nodes []map[string]any, runtimeID string) []map[string]any {
	runtimeID = strings.TrimSpace(runtimeID)
	if runtimeID == "" {
		return nodes
	}
	ownedWorkspaces := map[string]struct{}{runtimeID: struct{}{}}
	changed := true
	for changed {
		changed = false
		for _, node := range nodes {
			z := stringField(node, "z")
			if z == "" {
				continue
			}
			if _, owned := ownedWorkspaces[z]; !owned {
				continue
			}
			typ := stringField(node, "type")
			if !strings.HasPrefix(typ, "subflow:") {
				continue
			}
			id := strings.TrimSpace(strings.TrimPrefix(typ, "subflow:"))
			if id == "" {
				continue
			}
			if _, exists := ownedWorkspaces[id]; exists {
				continue
			}
			ownedWorkspaces[id] = struct{}{}
			changed = true
		}
	}
	out := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		typ := stringField(node, "type")
		id := stringField(node, "id")
		z := stringField(node, "z")
		if typ == "tab" {
			if id == runtimeID {
				out = append(out, node)
			}
			continue
		}
		if typ == "subflow" {
			if _, owned := ownedWorkspaces[id]; owned {
				out = append(out, node)
			}
			continue
		}
		if z != "" {
			if _, owned := ownedWorkspaces[z]; owned {
				out = append(out, node)
			}
			continue
		}
		out = append(out, node)
	}
	return out
}

func prepareNodeRedRuntimeNodes(nodes []map[string]any, runtimeID string) ([]map[string]any, map[string]struct{}, bool) {
	subflowDefinitions := make(map[string]struct{})
	ownedSubflows := make(map[string]struct{})
	for _, node := range nodes {
		if node == nil {
			continue
		}
		typ := stringField(node, "type")
		if typ == "subflow" {
			if id := stringField(node, "id"); id != "" {
				subflowDefinitions[id] = struct{}{}
				ownedSubflows[id] = struct{}{}
			}
			continue
		}
		if strings.HasPrefix(typ, "subflow:") {
			id := strings.TrimSpace(strings.TrimPrefix(typ, "subflow:"))
			if id != "" {
				ownedSubflows[id] = struct{}{}
			}
		}
	}
	runtimeNodes := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		typ := stringField(node, "type")
		if typ == "tab" {
			continue
		}
		if _, ok := node["z"]; ok {
			z := stringField(node, "z")
			if _, isSubflow := subflowDefinitions[z]; !isSubflow {
				node["z"] = runtimeID
			}
		}
		runtimeNodes = append(runtimeNodes, node)
	}
	return runtimeNodes, ownedSubflows, len(ownedSubflows) > 0
}

func nodeRedRootTab(runtimeID string, flow repo.Flow) map[string]any {
	return map[string]any{
		"id":       runtimeID,
		"type":     "tab",
		"label":    flow.Name,
		"disabled": false,
		"info":     flow.Description,
	}
}

func mergeNodeRedRuntimeNodes(current []map[string]any, rootTab map[string]any, runtimeNodes []map[string]any, ownedSubflows map[string]struct{}) []map[string]any {
	runtimeID := stringField(rootTab, "id")
	for id := range nodeRedRuntimeSubflows(current, runtimeID) {
		ownedSubflows[id] = struct{}{}
	}
	merged := make([]map[string]any, 0, len(current)+len(runtimeNodes)+1)
	for _, node := range current {
		if node == nil || isOwnedNodeRedRuntimeNode(node, runtimeID, ownedSubflows) {
			continue
		}
		merged = append(merged, node)
	}
	merged = append(merged, rootTab)
	merged = append(merged, runtimeNodes...)
	return merged
}

func nodeRedRuntimeSubflows(nodes []map[string]any, runtimeID string) map[string]struct{} {
	ids := make(map[string]struct{})
	for _, node := range nodes {
		if node == nil || stringField(node, "z") != runtimeID {
			continue
		}
		typ := stringField(node, "type")
		if !strings.HasPrefix(typ, "subflow:") {
			continue
		}
		id := strings.TrimSpace(strings.TrimPrefix(typ, "subflow:"))
		if id != "" {
			ids[id] = struct{}{}
		}
	}
	return ids
}

func isOwnedNodeRedRuntimeNode(node map[string]any, runtimeID string, ownedSubflows map[string]struct{}) bool {
	typ := stringField(node, "type")
	id := stringField(node, "id")
	if typ == "tab" && id == runtimeID {
		return true
	}
	z := stringField(node, "z")
	if z == runtimeID {
		return true
	}
	if typ == "subflow" {
		_, ok := ownedSubflows[id]
		return ok
	}
	if z != "" {
		_, ok := ownedSubflows[z]
		return ok
	}
	return false
}

func (s *Service) deployGlobalNodes(ctx context.Context, baseURL string, globalNodes []map[string]any, currentFlowRefs map[string]struct{}) error {
	if baseURL == "" || (len(globalNodes) == 0 && len(currentFlowRefs) == 0) {
		return nil
	}
	globalFlow, err := s.getNodeRed(ctx, baseURL, "/flow/global")
	if err != nil {
		return err
	}
	merged := s.mergeGlobalConfigs(ctx, baseURL, nodeRedGlobalConfigs(globalFlow), toAnySlice(globalNodes), currentFlowRefs)
	return s.updateNodeRedGlobalConfigs(ctx, baseURL, globalFlow, merged)
}

func nodeRedGlobalConfigs(globalFlow map[string]any) []any {
	configs := toAnySlice(globalFlow["configs"])
	if len(configs) == 0 {
		configs = toAnySlice(globalFlow["nodes"])
	}
	return configs
}

func (s *Service) updateNodeRedGlobalConfigs(ctx context.Context, baseURL string, globalFlow map[string]any, configs []any) error {
	if globalFlow == nil {
		globalFlow = map[string]any{}
	}
	// subflow 定义存储在 global flow 中，因此只替换 configs 并保留其余字段。
	globalFlow["id"] = "global"
	globalFlow["configs"] = configs
	delete(globalFlow, "nodes")
	_, err := s.postNodeRed(ctx, baseURL, http.MethodPut, "/flow/global", globalFlow)
	return err
}

func (s *Service) mergeGlobalConfigs(ctx context.Context, baseURL string, existing []any, incoming []any, currentFlowRefs map[string]struct{}) []any {
	installedTypes := s.fetchInstalledNodeTypes(ctx, baseURL)
	allFlows := []map[string]any(nil)
	if len(installedTypes) > 0 {
		allFlows, _ = s.getNodeRedFlowNodes(ctx, baseURL)
	}
	incomingMaps := filterMissingGlobalNodes(toMapSlice(incoming), installedTypes)
	incoming = toAnySlice(incomingMaps)
	incomingIDs := make(map[string]struct{}, len(incomingMaps))
	for _, cfg := range incomingMaps {
		if id := stringField(cfg, "id"); id != "" {
			incomingIDs[id] = struct{}{}
		}
	}
	merged := make([]any, 0, len(existing)+len(incoming))
	index := make(map[string]int, len(existing))
	for _, cfg := range existing {
		id := nodeID(cfg)
		nodeMap, _ := cfg.(map[string]any)
		if len(installedTypes) > 0 {
			if isMissingGlobalNode(nodeMap, installedTypes) {
				if _, referencedByCurrentFlow := currentFlowRefs[id]; referencedByCurrentFlow {
					continue
				}
				if isMissingOrphanGlobalNode(nodeMap, installedTypes, allFlows) {
					continue
				}
			}
		}
		if id != "" {
			index[id] = len(merged)
		}
		merged = append(merged, cfg)
	}
	for _, cfg := range incoming {
		id := nodeID(cfg)
		if id != "" {
			if pos, ok := index[id]; ok {
				merged[pos] = cfg
				continue
			}
		}
		merged = append(merged, cfg)
	}
	return merged
}

func (s *Service) nodeRedFlowExists(ctx context.Context, baseURL, runtimeID string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(baseURL, "/")+"/flow/"+urlPath(runtimeID), nil)
	if err != nil {
		return false, err
	}
	setNodeRedInternalAuth(req)
	resp, err := s.http.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return false, fmt.Errorf("node-red GET /flow/%s: status=%d body=%s", runtimeID, resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return true, nil
}

func (s *Service) deleteNodeRedRuntimeFlow(ctx context.Context, flowType int, runtimeID string) error {
	runtimeID = strings.TrimSpace(runtimeID)
	if runtimeID == "" {
		return nil
	}
	baseURL := s.runtimeBaseURL(flowType)
	if baseURL == "" {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, strings.TrimRight(baseURL, "/")+"/flow/"+urlPath(runtimeID), nil)
	if err != nil {
		return err
	}
	setNodeRedInternalAuth(req)
	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil
	}
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Errorf("node-red DELETE /flow/%s: status=%d body=%s", runtimeID, resp.StatusCode, strings.TrimSpace(string(respBody)))
}

func (s *Service) runtimeBaseURL(flowType int) string {
	if flowType == 2 {
		return s.eventFlowURL
	}
	return s.sourceFlowURL
}

func (s *Service) runtimeRevision(ctx context.Context, flowType int) (string, error) {
	baseURL := s.runtimeBaseURL(flowType)
	if baseURL == "" {
		return "", nil
	}
	out, err := s.getNodeRed(ctx, baseURL, "/flows")
	if err != nil {
		return "", err
	}
	return stringField(out, "rev"), nil
}

func (s *Service) getNodeRedFlowNodes(ctx context.Context, baseURL string) ([]map[string]any, error) {
	out, err := s.getNodeRed(ctx, baseURL, "/flows")
	if err != nil {
		return nil, err
	}
	return runtimeNodesFromAny(out["flows"])
}

func runtimeNodesFromAny(value any) ([]map[string]any, error) {
	if value == nil {
		return []map[string]any{}, nil
	}
	if nodes, ok := value.([]map[string]any); ok {
		return nodes, nil
	}
	items, ok := value.([]any)
	if !ok {
		return nil, errors.New("node-red returned invalid flows payload")
	}
	nodes := make([]map[string]any, 0, len(items))
	for _, item := range items {
		node, ok := item.(map[string]any)
		if !ok {
			return nil, errors.New("node-red returned invalid flow node")
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (s *Service) postNodeRedFlows(ctx context.Context, baseURL string, nodes []map[string]any) error {
	nodes = normalizeSubflowNodes(nodes)
	body, err := json.Marshal(map[string]any{"flows": nodes})
	if err != nil {
		return err
	}
	headers := map[string]string{
		"Content-Type":             "application/json",
		"Node-RED-API-Version":     "v2",
		"Node-RED-Deployment-Type": "flows",
	}
	_, err = s.doNodeRedBytes(ctx, baseURL, http.MethodPost, "/flows", body, headers)
	return err
}

func (s *Service) postNodeRed(ctx context.Context, baseURL, method, path string, payload any) (map[string]any, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	respBody, err := s.doNodeRedBytes(ctx, baseURL, method, path, body, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if len(bytes.TrimSpace(respBody)) > 0 {
		_ = json.Unmarshal(respBody, &out)
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, nil
}

func (s *Service) getNodeRed(ctx context.Context, baseURL, path string) (map[string]any, error) {
	respBody, err := s.doNodeRedBytes(ctx, baseURL, http.MethodGet, path, nil, map[string]string{"Node-RED-API-Version": "v2"})
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if len(bytes.TrimSpace(respBody)) > 0 {
		_ = json.Unmarshal(respBody, &out)
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, nil
}

func (s *Service) doNodeRedBytes(ctx context.Context, baseURL, method, path string, body []byte, headers map[string]string) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(baseURL, "/")+path, reader)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	setNodeRedInternalAuth(req)
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("node-red %s %s: status=%d body=%s", method, path, resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return respBody, nil
}

func setNodeRedInternalAuth(req *http.Request) {
	if token := secrets.InternalToken("NODERED_INTERNAL_TOKEN"); token != "" {
		req.Header.Set(nodeRedInternalTokenHeader, token)
	}
}

func nodeRedCredentialsFromNodes(nodes []map[string]any) map[string]map[string]string {
	credentials := make(map[string]map[string]string)
	for _, node := range nodes {
		if node == nil {
			continue
		}
		id := stringField(node, "id")
		if id == "" {
			continue
		}
		fields := nodeRedCredentialFields(node["credentials"])
		if len(fields) > 0 {
			credentials[id] = fields
		}
	}
	return credentials
}

func nodeRedCredentialFields(raw any) map[string]string {
	switch v := raw.(type) {
	case map[string]any:
		fields := make(map[string]string, len(v))
		for key, value := range v {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}
			if text, ok := nodeRedCredentialValue(value); ok {
				fields[key] = text
			}
		}
		return fields
	case map[string]string:
		fields := make(map[string]string, len(v))
		for key, value := range v {
			key = strings.TrimSpace(key)
			if key != "" {
				fields[key] = value
			}
		}
		return fields
	default:
		return nil
	}
}

func nodeRedCredentialValue(value any) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case int:
		return strconv.Itoa(v), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case bool:
		if v {
			return "true", true
		}
		return "false", true
	default:
		return "", false
	}
}

func removeWireReferences(nodes []map[string]any, removedIDs map[string]struct{}) {
	if len(nodes) == 0 || len(removedIDs) == 0 {
		return
	}
	for _, node := range nodes {
		wires, ok := node["wires"].([]any)
		if !ok {
			continue
		}
		filteredWires := make([]any, 0, len(wires))
		for _, group := range wires {
			targets, ok := group.([]any)
			if !ok {
				filteredWires = append(filteredWires, group)
				continue
			}
			filteredTargets := make([]any, 0, len(targets))
			for _, target := range targets {
				if _, removed := removedIDs[strings.TrimSpace(fmt.Sprint(target))]; removed {
					continue
				}
				filteredTargets = append(filteredTargets, target)
			}
			filteredWires = append(filteredWires, filteredTargets)
		}
		node["wires"] = filteredWires
	}
}

func referencedNodeIDs(nodes []map[string]any) map[string]struct{} {
	refs := make(map[string]struct{})
	for _, node := range nodes {
		collectReferences(refs, node, "id")
	}
	return refs
}

func nodeReferencesID(node map[string]any, id string) bool {
	if node == nil || id == "" {
		return false
	}
	return valueReferencesID(node, id, "id")
}

func collectReferences(refs map[string]struct{}, value any, skippedKey string) {
	switch v := value.(type) {
	case map[string]any:
		for key, item := range v {
			if key == skippedKey {
				continue
			}
			collectReferences(refs, item, "")
		}
	case []any:
		for _, item := range v {
			collectReferences(refs, item, "")
		}
	case []map[string]any:
		for _, item := range v {
			collectReferences(refs, item, "")
		}
	case string:
		s := strings.TrimSpace(v)
		if s != "" {
			refs[s] = struct{}{}
		}
	}
}

func valueReferencesID(value any, id string, skippedKey string) bool {
	switch v := value.(type) {
	case map[string]any:
		for key, item := range v {
			if key == skippedKey {
				continue
			}
			if valueReferencesID(item, id, "") {
				return true
			}
		}
	case []any:
		for _, item := range v {
			if valueReferencesID(item, id, "") {
				return true
			}
		}
	case []map[string]any:
		for _, item := range v {
			if valueReferencesID(item, id, "") {
				return true
			}
		}
	case string:
		return strings.TrimSpace(v) == id
	}
	return false
}

func stringField(node map[string]any, key string) string {
	if node == nil {
		return ""
	}
	value, ok := node[key]
	if !ok || value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func nodeID(node any) string {
	if m, ok := node.(map[string]any); ok {
		return stringField(m, "id")
	}
	return ""
}

func toMapSlice(value any) []map[string]any {
	switch v := value.(type) {
	case []map[string]any:
		return v
	case []any:
		res := make([]map[string]any, 0, len(v))
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				res = append(res, m)
			}
		}
		return res
	default:
		return nil
	}
}

func toAnySlice(value any) []any {
	switch v := value.(type) {
	case []any:
		return v
	case []map[string]any:
		res := make([]any, 0, len(v))
		for _, item := range v {
			res = append(res, item)
		}
		return res
	default:
		return nil
	}
}

func urlPath(value string) string {
	return strings.ReplaceAll(strings.TrimSpace(value), "/", "%2F")
}

func asString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	default:
		return ""
	}
}

func sortedStrings(values map[string]struct{}) []string {
	out := make([]string, 0, len(values))
	for value := range values {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
