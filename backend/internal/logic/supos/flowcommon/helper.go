package flowcommon

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"backend/internal/common"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"

	noderedclient "backend/share/clients/nodered"

	"gitee.com/unitedrhino/share/errors"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	FlowStatusDraft   = "DRAFT"
	FlowStatusPending = "PENDING"
	FlowStatusRunning = "RUNNING"
)

// FlowRepo declares the repository behaviour shared by Node-RED flows.
type FlowRepo interface {
	FindOne(ctx context.Context, id int64) (*relationDB.NoderedFlow, error)
	Insert(ctx context.Context, data *relationDB.NoderedFlow) error
	Update(ctx context.Context, data *relationDB.NoderedFlow) error
	ReplaceModels(ctx context.Context, parentID int64, aliases []string) error
}

// FlowCopyInput defines the required fields when copying a flow.
type FlowCopyInput struct {
	FlowName    string
	Description string
	Template    string
	GroupId     *int64
	Creator     string
}

// CopyFlow clones the given flow and returns the created record.
func CopyFlow(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	repo FlowRepo,
	sourceID int64,
	input FlowCopyInput,
	client *noderedclient.Client,
) (*relationDB.NoderedFlow, error) {
	src, err := repo.FindOne(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	if src == nil {
		return nil, errors.NotFind.WithMsg("nodered.flow.not.exist")
	}

	dst := &relationDB.NoderedFlow{
		ID:          common.NextId(),
		FlowName:    strings.TrimSpace(input.FlowName),
		Description: strings.TrimSpace(input.Description),
		Template:    strings.TrimSpace(input.Template),
		FlowStatus:  FlowStatusDraft,
		GroupId:     input.GroupId,
		Creator:     strings.TrimSpace(input.Creator),
	}

	sourceJSON, _, err := ResolveNodesJSON(ctx, client, "", src)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(sourceJSON) != "" {
		newJSON, err := regenerateNodeIDs(sourceJSON)
		if err != nil {
			return nil, err
		}
		dst.FlowData = newJSON
	}

	if err := repo.Insert(ctx, dst); err != nil {
		return nil, err
	}
	return dst, nil
}

// DeployFlow pushes the flow definition to Node-RED and persists the latest state.
func DeployFlow(
	ctx context.Context,
	repo FlowRepo,
	entityID int64,
	overrideJSON string,
	client *noderedclient.Client,
	aliasExtractor func([]map[string]any) []string,
) (string, error) {
	if client == nil {
		return "", errors.System.WithMsg("nodered.flow.not.exist")
	}

	rec, err := repo.FindOne(ctx, entityID)
	if err != nil {
		return "", err
	}
	if rec == nil {
		return "", errors.NotFind.WithMsg("nodered.flow.not.exist")
	}

	resolvedJSON, _, err := ResolveNodesJSON(ctx, client, overrideJSON, rec)
	if err != nil {
		return "", err
	}
	resolvedJSON = strings.TrimSpace(resolvedJSON)
	if resolvedJSON == "" {
		return "", errors.Parameter.WithMsg("nodered.flowId.empty")
	}

	var rawNodes []map[string]any
	if err := json.Unmarshal([]byte(resolvedJSON), &rawNodes); err != nil {
		return "", errors.Parameter.WithMsg("nodered.invalid.parameter")
	}
	if err := ValidateNoMissingNodeTypes(ctx, client, rawNodes); err != nil {
		return "", err
	}

	flowNodes, globalNodes := splitGlobalNodes(rawNodes)

	// mqtt的全局节点先部署
	if len(globalNodes) > 0 {
		currentFlowRefs := ReferencedNodeIDs(flowNodes)
		mergedGlobal := mergeGlobalNodes(ctx, client, globalNodes, currentFlowRefs)
		globalBody := map[string]any{
			"id":      "global",
			"configs": toInterfaceSlice(mergedGlobal),
		}
		var gout map[string]any
		code, body, errs := client.DoJSON(ctx, "PUT", "/flow/global", globalBody, &gout)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("update global flow failed: code=%d err=%v body=%s", code, errs, string(body))
			return "", errors.System.WithMsg("error.sys.systemError").AddDetailf("node-red update global failed: code=%d err=%v body=%s", code, errs, string(body))
		}
	}

	flowID := strings.TrimSpace(rec.FlowID)
	// create flow if absent
	if flowID == "" {
		req := map[string]any{
			"id":       "",
			"nodes":    []any{},
			"disabled": false,
			"label":    rec.FlowName,
			"info":     rec.Description,
		}
		var out map[string]any
		code, body, errs := client.DoJSON(ctx, "POST", "/flow", req, &out)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("create flow failed: code=%d err=%v body=%s", code, errs, string(body))
			return "", errors.System.WithMsg("error.sys.systemError").AddDetailf("node-red create flow failed: code=%d err=%v body=%s", code, errs, string(body))
		}
		if id, ok := out["id"].(string); ok && strings.TrimSpace(id) != "" {
			flowID = id
		} else {
			return "", errors.System.WithMsg("error.sys.systemError").AddDetail("node-red create flow returned empty id")
		}
	}

	setZ(flowNodes, flowID)

	flowBody := map[string]any{
		"id":       flowID,
		"nodes":    toInterfaceSlice(flowNodes),
		"disabled": false,
		"label":    rec.FlowName,
		"info":     rec.Description,
	}
	var upd map[string]any
	code, body, errs := client.DoJSON(ctx, "PUT", "/flow/"+flowID, flowBody, &upd)
	if len(errs) > 0 || (code != 200 && code != 204) {
		logx.WithContext(ctx).Errorf("update flow failed: code=%d err=%v body=%s", code, errs, string(body))
		return "", errors.System.WithMsg("error.sys.systemError").AddDetailf("node-red update flow failed: code=%d err=%v body=%s", code, errs, string(body))
	}

	rec.FlowID = flowID
	rec.FlowStatus = FlowStatusRunning
	rec.FlowData = ""
	if err := repo.Update(ctx, rec); err != nil {
		return "", err
	}

	return flowID, nil
}

// ResolveNodesJSON resolves the nodes JSON string that should be used for deploy/save operations.
func ResolveNodesJSON(ctx context.Context, client *noderedclient.Client, override string, entity *relationDB.NoderedFlow) (string, string, error) {
	override = strings.TrimSpace(override)
	if override != "" {
		return override, "client", nil
	}
	raw := strings.TrimSpace(entity.FlowData)
	if raw != "" {
		return raw, "draft", nil
	}
	// fetch from node-red runtime
	if client != nil && strings.TrimSpace(entity.FlowID) != "" {
		var out map[string]any
		code, body, errs := client.GetFlowNodesV1(ctx, entity.FlowID, &out)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("fetch nodes from node-red failed: code=%d err=%v body=%s", code, errs, string(body))
			return "", "", errors.System.WithMsg("nodered.flow.not.exist")
		}
		if nodes, ok := out["nodes"].([]any); ok {
			js, err := json.Marshal(nodes)
			if err != nil {
				return "", "", errors.System.WithMsg(err.Error())
			}
			return string(js), "nodered", nil
		}
	}
	return "", "", nil
}

// GenerateNodeID generates a random Node-RED node id (16 hex chars).
func GenerateNodeID() string {
	u := strings.ReplaceAll(uuid.NewString(), "-", "")
	if len(u) > 16 {
		return u[:16]
	}
	return u
}

func regenerateNodeIDs(jsonStr string) (string, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return jsonStr, nil
	}
	var nodes []map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &nodes); err != nil {
		return "", errors.Parameter.WithMsg("nodered.invalid.parameter")
	}
	result := jsonStr
	for _, node := range nodes {
		z, _ := node["z"].(string)
		if strings.TrimSpace(z) == "" {
			continue
		}
		oldID, _ := node["id"].(string)
		if strings.TrimSpace(oldID) == "" {
			continue
		}
		newID := GenerateNodeID()
		result = strings.ReplaceAll(result, oldID, newID)
	}
	return result, nil
}

func splitGlobalNodes(nodes []map[string]any) (flowNodes []map[string]any, globalNodes []map[string]any) {
	flowNodes = make([]map[string]any, 0, len(nodes))
	globalNodes = make([]map[string]any, 0)
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if _, ok := node["z"]; ok {
			flowNodes = append(flowNodes, node)
			continue
		}
		t, _ := node["type"].(string)
		if strings.TrimSpace(t) == "tab" {
			flowNodes = append(flowNodes, node)
			continue
		}
		globalNodes = append(globalNodes, node)
	}
	return
}

func setZ(nodes []map[string]any, flowID string) {
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if _, ok := node["z"]; ok {
			node["z"] = flowID
		}
	}
}

func toInterfaceSlice(nodes []map[string]any) []any {
	out := make([]any, 0, len(nodes))
	for _, n := range nodes {
		if n != nil {
			out = append(out, n)
		}
	}
	return out
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

// ExtractAliases parses possible UNS aliases from Node-RED node definitions.
func ExtractAliases(nodes []map[string]any) []string {
	aliasSet := make(map[string]struct{})
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if val, ok := node["selectedModelAlias"]; ok {
			if alias := strings.TrimSpace(fmt.Sprint(val)); alias != "" {
				aliasSet[alias] = struct{}{}
			}
		}
	}
	aliases := make([]string, 0, len(aliasSet))
	for alias := range aliasSet {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	return aliases
}

func mergeGlobalNodes(ctx context.Context, client *noderedclient.Client, incoming []map[string]any, currentFlowRefs map[string]struct{}) []map[string]any {
	existing := fetchGlobalNodes(ctx, client)
	installedTypes := FetchInstalledNodeTypes(ctx, client)
	if len(installedTypes) == 0 {
		installedTypes = nil
	}
	incoming = filterMissingGlobalNodes(ctx, incoming, installedTypes)
	if len(existing) == 0 {
		return incoming
	}
	allFlows := fetchAllFlowNodes(ctx, client)

	incomingByID := make(map[string]map[string]any)
	incomingNoID := make([]map[string]any, 0)
	for _, node := range incoming {
		if node == nil {
			continue
		}
		id := strings.TrimSpace(fmt.Sprint(node["id"]))
		if id == "" {
			incomingNoID = append(incomingNoID, node)
			continue
		}
		incomingByID[id] = node
	}

	merged := make([]map[string]any, 0, len(existing)+len(incoming))
	for _, node := range existing {
		if node == nil {
			continue
		}
		id := strings.TrimSpace(fmt.Sprint(node["id"]))
		if id == "" {
			merged = append(merged, node)
			continue
		}
		if updated, ok := incomingByID[id]; ok {
			merged = append(merged, updated)
			delete(incomingByID, id)
		} else {
			if installedTypes != nil && isMissingGlobalNode(node, installedTypes) {
				if _, referencedByCurrentFlow := currentFlowRefs[id]; referencedByCurrentFlow {
					continue
				}
				if isMissingOrphanGlobalNode(node, installedTypes, allFlows) {
					continue
				}
			}
			merged = append(merged, node)
		}
	}
	for _, node := range incomingNoID {
		merged = append(merged, node)
	}
	for _, node := range incomingByID {
		merged = append(merged, node)
	}
	return merged
}

// ValidateNoMissingNodeTypes rejects nodes whose type is not installed.
func ValidateNoMissingNodeTypes(ctx context.Context, client *noderedclient.Client, nodes []map[string]any) error {
	if client == nil || len(nodes) == 0 {
		return nil
	}
	installedTypes := FetchInstalledNodeTypes(ctx, client)
	if len(installedTypes) == 0 {
		return nil
	}
	missingTypes := MissingNodeTypes(nodes, installedTypes)
	if len(missingTypes) == 0 {
		return nil
	}
	return errors.Parameter.WithMsg("missing Node-RED node types: " + strings.Join(missingTypes, ", "))
}

func MissingNodeTypes(nodes []map[string]any, installedTypes map[string]bool) []string {
	if len(nodes) == 0 || len(installedTypes) == 0 {
		return nil
	}
	missing := make(map[string]struct{})
	for _, node := range nodes {
		if !isMissingNodeType(node, installedTypes) {
			continue
		}
		typ := strings.TrimSpace(fmt.Sprint(node["type"]))
		if typ != "" {
			missing[typ] = struct{}{}
		}
	}
	if len(missing) == 0 {
		return nil
	}
	types := make([]string, 0, len(missing))
	for typ := range missing {
		types = append(types, typ)
	}
	sort.Strings(types)
	return types
}

// RemoveMissingRuntimeNodes removes nodes whose type is not installed.
//
// A missing node type stops the entire Node-RED runtime, even when unrelated
// tabs do not reference it. Removing the unknown entries keeps unrelated flows
// runnable while preserving nodes with installed types.
func RemoveMissingRuntimeNodes(ctx context.Context, client *noderedclient.Client) (int, error) {
	if client == nil {
		return 0, nil
	}
	installedTypes := FetchInstalledNodeTypes(ctx, client)
	if len(installedTypes) == 0 {
		return 0, nil
	}
	removedGlobal, err := removeMissingGlobalNodes(ctx, client, installedTypes)
	if err != nil {
		return 0, err
	}
	removedFlows, err := removeMissingFlowNodes(ctx, client, installedTypes)
	if err != nil {
		return removedGlobal, err
	}
	return removedGlobal + removedFlows, nil
}

// RuntimeMissingNode describes an unknown Node-RED node found anywhere in the runtime.
type RuntimeMissingNode struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name,omitempty"`
	Scope     string `json:"scope"`
	FlowID    string `json:"flowId,omitempty"`
	FlowLabel string `json:"flowLabel,omitempty"`
	Users     int    `json:"users,omitempty"`
}

// MissingNodeDeleteTarget identifies a missing node to remove from the runtime.
type MissingNodeDeleteTarget struct {
	ID     string
	FlowID string
	Scope  string
}

const (
	MissingNodeScopeGlobalConfig = "globalConfig"
	MissingNodeScopeFlowNode     = "flowNode"
	MissingNodeScopeFlowConfig   = "flowConfig"
)

// ListMissingRuntimeNodes returns all runtime nodes whose type is not installed.
func ListMissingRuntimeNodes(ctx context.Context, client *noderedclient.Client) ([]RuntimeMissingNode, error) {
	if client == nil {
		return nil, nil
	}
	installedTypes := FetchInstalledNodeTypes(ctx, client)
	if len(installedTypes) == 0 {
		return nil, nil
	}

	missing := make([]RuntimeMissingNode, 0)
	for _, node := range fetchGlobalNodes(ctx, client) {
		if isMissingGlobalNode(node, installedTypes) {
			missing = append(missing, runtimeMissingNode(node, MissingNodeScopeGlobalConfig, "global", "On all flows"))
		}
	}

	for _, tab := range fetchFlowTabs(ctx, client) {
		flowID := strings.TrimSpace(fmt.Sprint(tab["id"]))
		if flowID == "" {
			continue
		}
		var out map[string]any
		code, body, errs := client.GetFlowNodesV1(ctx, flowID, &out)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("fetch flow missing nodes failed: flow=%s code=%d err=%v body=%s", flowID, code, errs, string(body))
			return nil, errors.System.WithMsg("error.sys.systemError").AddDetailf("node-red fetch flow failed: flow=%s code=%d err=%v body=%s", flowID, code, errs, string(body))
		}
		flowLabel := strings.TrimSpace(fmt.Sprint(out["label"]))
		if flowLabel == "" {
			flowLabel = strings.TrimSpace(fmt.Sprint(tab["label"]))
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

// DeleteMissingRuntimeNode removes one unknown node by id and location.
func DeleteMissingRuntimeNode(ctx context.Context, client *noderedclient.Client, target MissingNodeDeleteTarget) (int, error) {
	if client == nil {
		return 0, nil
	}
	target.ID = strings.TrimSpace(target.ID)
	target.FlowID = strings.TrimSpace(target.FlowID)
	target.Scope = strings.TrimSpace(target.Scope)
	if target.ID == "" {
		return 0, errors.Parameter.WithMsg("nodered.nodeId.empty")
	}
	installedTypes := FetchInstalledNodeTypes(ctx, client)
	if len(installedTypes) == 0 {
		return 0, nil
	}
	if target.FlowID == "" && target.Scope != MissingNodeScopeGlobalConfig {
		resolved, err := resolveMissingNodeTarget(ctx, client, installedTypes, target.ID)
		if err != nil {
			return 0, err
		}
		target = resolved
	}
	if target.Scope == MissingNodeScopeGlobalConfig || target.FlowID == "global" {
		return deleteMissingGlobalNode(ctx, client, installedTypes, target.ID)
	}
	if target.FlowID == "" {
		return 0, errors.Parameter.WithMsg("nodered.flowId.empty")
	}
	return deleteMissingFlowNode(ctx, client, installedTypes, target)
}

// RemoveMissingConfigNodes is kept for older callers; it now also removes
// normal flow nodes with unknown types because those stop the whole runtime too.
func RemoveMissingConfigNodes(ctx context.Context, client *noderedclient.Client) (int, error) {
	return RemoveMissingRuntimeNodes(ctx, client)
}

// RemoveMissingGlobalNodes is kept for older callers; new code should use
// RemoveMissingRuntimeNodes so tab-scoped config and normal nodes are handled as well.
func RemoveMissingGlobalNodes(ctx context.Context, client *noderedclient.Client) (int, error) {
	if client == nil {
		return 0, nil
	}
	installedTypes := FetchInstalledNodeTypes(ctx, client)
	if len(installedTypes) == 0 {
		return 0, nil
	}
	return removeMissingGlobalNodes(ctx, client, installedTypes)
}

func removeMissingGlobalNodes(ctx context.Context, client *noderedclient.Client, installedTypes map[string]bool) (int, error) {
	existing := fetchGlobalNodes(ctx, client)
	if len(existing) == 0 {
		return 0, nil
	}
	filtered := make([]map[string]any, 0, len(existing))
	removed := make([]string, 0)
	for _, node := range existing {
		if isMissingGlobalNode(node, installedTypes) {
			removed = append(removed, fmt.Sprintf("%s(%s)", strings.TrimSpace(fmt.Sprint(node["id"])), strings.TrimSpace(fmt.Sprint(node["type"]))))
			continue
		}
		filtered = append(filtered, node)
	}
	if len(removed) == 0 {
		return 0, nil
	}

	body := map[string]any{
		"id":      "global",
		"configs": toInterfaceSlice(filtered),
	}
	var out map[string]any
	code, respBody, errs := client.DoJSON(ctx, "PUT", "/flow/global", body, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		logx.WithContext(ctx).Errorf("remove missing global nodes failed: removed=%v code=%d err=%v body=%s", removed, code, errs, string(respBody))
		return 0, errors.System.WithMsg("error.sys.systemError").AddDetailf("node-red remove missing global failed: code=%d err=%v body=%s", code, errs, string(respBody))
	}
	logx.WithContext(ctx).Infof("removed missing node-red global configs: %v", removed)
	return len(removed), nil
}

func deleteMissingGlobalNode(ctx context.Context, client *noderedclient.Client, installedTypes map[string]bool, nodeID string) (int, error) {
	existing := fetchGlobalNodes(ctx, client)
	if len(existing) == 0 {
		return 0, nil
	}
	filtered := make([]map[string]any, 0, len(existing))
	removed := false
	for _, node := range existing {
		if strings.TrimSpace(fmt.Sprint(node["id"])) == nodeID && isMissingGlobalNode(node, installedTypes) {
			removed = true
			continue
		}
		filtered = append(filtered, node)
	}
	if !removed {
		return 0, nil
	}
	body := map[string]any{
		"id":      "global",
		"configs": toInterfaceSlice(filtered),
	}
	var out map[string]any
	code, respBody, errs := client.DoJSON(ctx, "PUT", "/flow/global", body, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		logx.WithContext(ctx).Errorf("delete missing global node failed: node=%s code=%d err=%v body=%s", nodeID, code, errs, string(respBody))
		return 0, errors.System.WithMsg("error.sys.systemError").AddDetailf("node-red delete missing global failed: code=%d err=%v body=%s", code, errs, string(respBody))
	}
	return 1, nil
}

func removeMissingFlowNodes(ctx context.Context, client *noderedclient.Client, installedTypes map[string]bool) (int, error) {
	tabs := fetchFlowTabs(ctx, client)
	removedCount := 0
	for _, tab := range tabs {
		flowID := strings.TrimSpace(fmt.Sprint(tab["id"]))
		if flowID == "" {
			continue
		}
		var out map[string]any
		code, body, errs := client.GetFlowNodesV1(ctx, flowID, &out)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("fetch flow configs failed: flow=%s code=%d err=%v body=%s", flowID, code, errs, string(body))
			continue
		}
		nodes := toMapSlice(out["nodes"])
		configs := toMapSlice(out["configs"])
		filteredNodes, removedNodes, removedIDs := filterMissingRuntimeNodes(nodes, installedTypes)
		filteredConfigs, removedConfigs, removedConfigIDs := filterMissingRuntimeNodes(configs, installedTypes)
		if len(removedNodes) == 0 && len(removedConfigs) == 0 {
			continue
		}
		for id := range removedConfigIDs {
			removedIDs[id] = struct{}{}
		}
		removeWireReferences(filteredNodes, removedIDs)
		reqBody := map[string]any{
			"id":       flowID,
			"label":    out["label"],
			"disabled": out["disabled"],
			"info":     out["info"],
			"nodes":    toInterfaceSlice(filteredNodes),
			"configs":  toInterfaceSlice(filteredConfigs),
		}
		var updated map[string]any
		code, respBody, errs := client.DoJSON(ctx, "PUT", "/flow/"+flowID, reqBody, &updated)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("remove missing flow nodes failed: flow=%s removedNodes=%v removedConfigs=%v code=%d err=%v body=%s", flowID, removedNodes, removedConfigs, code, errs, string(respBody))
			return removedCount, errors.System.WithMsg("error.sys.systemError").AddDetailf("node-red remove missing flow nodes failed: flow=%s code=%d err=%v body=%s", flowID, code, errs, string(respBody))
		}
		removedCount += len(removedNodes) + len(removedConfigs)
		logx.WithContext(ctx).Infof("removed missing node-red flow nodes: flow=%s nodes=%v configs=%v", flowID, removedNodes, removedConfigs)
	}
	return removedCount, nil
}

func deleteMissingFlowNode(ctx context.Context, client *noderedclient.Client, installedTypes map[string]bool, target MissingNodeDeleteTarget) (int, error) {
	var out map[string]any
	code, body, errs := client.GetFlowNodesV1(ctx, target.FlowID, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		logx.WithContext(ctx).Errorf("fetch flow for missing node delete failed: flow=%s code=%d err=%v body=%s", target.FlowID, code, errs, string(body))
		return 0, errors.System.WithMsg("error.sys.systemError").AddDetailf("node-red fetch flow failed: flow=%s code=%d err=%v body=%s", target.FlowID, code, errs, string(body))
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
	reqBody := map[string]any{
		"id":       target.FlowID,
		"label":    out["label"],
		"disabled": out["disabled"],
		"info":     out["info"],
		"nodes":    toInterfaceSlice(filteredNodes),
		"configs":  toInterfaceSlice(filteredConfigs),
	}
	var updated map[string]any
	code, respBody, errs := client.DoJSON(ctx, "PUT", "/flow/"+target.FlowID, reqBody, &updated)
	if len(errs) > 0 || (code != 200 && code != 204) {
		logx.WithContext(ctx).Errorf("delete missing flow node failed: flow=%s node=%s code=%d err=%v body=%s", target.FlowID, target.ID, code, errs, string(respBody))
		return 0, errors.System.WithMsg("error.sys.systemError").AddDetailf("node-red delete missing flow node failed: flow=%s code=%d err=%v body=%s", target.FlowID, code, errs, string(respBody))
	}
	return len(removedNodeIDs) + len(removedConfigIDs), nil
}

func filterMissingRuntimeNodes(nodes []map[string]any, installedTypes map[string]bool) ([]map[string]any, []string, map[string]struct{}) {
	filtered := make([]map[string]any, 0, len(nodes))
	removed := make([]string, 0)
	removedIDs := make(map[string]struct{})
	for _, node := range nodes {
		if isMissingNodeType(node, installedTypes) {
			id := strings.TrimSpace(fmt.Sprint(node["id"]))
			if id != "" {
				removedIDs[id] = struct{}{}
			}
			removed = append(removed, fmt.Sprintf("%s(%s)", id, strings.TrimSpace(fmt.Sprint(node["type"]))))
			continue
		}
		filtered = append(filtered, node)
	}
	return filtered, removed, removedIDs
}

func removeMissingNodeByID(nodes []map[string]any, installedTypes map[string]bool, nodeID string) ([]map[string]any, map[string]struct{}) {
	filtered := make([]map[string]any, 0, len(nodes))
	removedIDs := make(map[string]struct{})
	for _, node := range nodes {
		if strings.TrimSpace(fmt.Sprint(node["id"])) == nodeID && isMissingNodeType(node, installedTypes) {
			removedIDs[nodeID] = struct{}{}
			continue
		}
		filtered = append(filtered, node)
	}
	return filtered, removedIDs
}

func resolveMissingNodeTarget(ctx context.Context, client *noderedclient.Client, installedTypes map[string]bool, nodeID string) (MissingNodeDeleteTarget, error) {
	matches := make([]MissingNodeDeleteTarget, 0, 1)
	for _, node := range fetchGlobalNodes(ctx, client) {
		if strings.TrimSpace(fmt.Sprint(node["id"])) == nodeID && isMissingGlobalNode(node, installedTypes) {
			matches = append(matches, MissingNodeDeleteTarget{ID: nodeID, FlowID: "global", Scope: MissingNodeScopeGlobalConfig})
		}
	}
	for _, tab := range fetchFlowTabs(ctx, client) {
		flowID := strings.TrimSpace(fmt.Sprint(tab["id"]))
		if flowID == "" {
			continue
		}
		var out map[string]any
		code, body, errs := client.GetFlowNodesV1(ctx, flowID, &out)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("resolve missing node target failed: flow=%s code=%d err=%v body=%s", flowID, code, errs, string(body))
			continue
		}
		for _, node := range toMapSlice(out["nodes"]) {
			if strings.TrimSpace(fmt.Sprint(node["id"])) == nodeID && isMissingNodeType(node, installedTypes) {
				matches = append(matches, MissingNodeDeleteTarget{ID: nodeID, FlowID: flowID, Scope: MissingNodeScopeFlowNode})
			}
		}
		for _, node := range toMapSlice(out["configs"]) {
			if strings.TrimSpace(fmt.Sprint(node["id"])) == nodeID && isMissingNodeType(node, installedTypes) {
				matches = append(matches, MissingNodeDeleteTarget{ID: nodeID, FlowID: flowID, Scope: MissingNodeScopeFlowConfig})
			}
		}
	}
	if len(matches) == 0 {
		return MissingNodeDeleteTarget{}, errors.NotFind.WithMsg("nodered.node.not.exist")
	}
	if len(matches) > 1 {
		return MissingNodeDeleteTarget{}, errors.Parameter.WithMsg("nodered.node.target.ambiguous")
	}
	return matches[0], nil
}

func runtimeMissingNode(node map[string]any, scope, flowID, flowLabel string) RuntimeMissingNode {
	id := strings.TrimSpace(fmt.Sprint(node["id"]))
	return RuntimeMissingNode{
		ID:        id,
		Type:      strings.TrimSpace(fmt.Sprint(node["type"])),
		Name:      missingNodeDisplayName(node, id),
		Scope:     scope,
		FlowID:    flowID,
		FlowLabel: flowLabel,
		Users:     len(toAnySlice(node["_users"])),
	}
}

func missingNodeDisplayName(node map[string]any, fallback string) string {
	for _, key := range []string{"name", "label", "topic"} {
		name := strings.TrimSpace(fmt.Sprint(node[key]))
		if name != "" && name != "<nil>" {
			return name
		}
	}
	return fallback
}

func filterMissingGlobalNodes(ctx context.Context, nodes []map[string]any, installedTypes map[string]bool) []map[string]any {
	if len(nodes) == 0 || len(installedTypes) == 0 {
		return nodes
	}
	filtered := make([]map[string]any, 0, len(nodes))
	removed := make([]string, 0)
	for _, node := range nodes {
		if isMissingGlobalNode(node, installedTypes) {
			removed = append(removed, fmt.Sprintf("%s(%s)", strings.TrimSpace(fmt.Sprint(node["id"])), strings.TrimSpace(fmt.Sprint(node["type"]))))
			continue
		}
		filtered = append(filtered, node)
	}
	if len(removed) > 0 {
		logx.WithContext(ctx).Infof("skip missing node-red global configs from deploy payload: %v", removed)
	}
	return filtered
}

// FetchInstalledNodeTypes returns the node types currently registered by Node-RED.
func FetchInstalledNodeTypes(ctx context.Context, client *noderedclient.Client) map[string]bool {
	installed := make(map[string]bool)
	if client == nil {
		return installed
	}
	code, body, errs := client.DoBytes(ctx, "GET", "/nodes", nil)
	if len(errs) > 0 || (code != 200 && code != 204) {
		logx.WithContext(ctx).Errorf("fetch node-red installed nodes failed: code=%d err=%v body=%s", code, errs, string(body))
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
		collectInstalledNodeTypesFromNodeSets(installed, out)
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

func collectInstalledNodeTypesFromNodeSets(installed map[string]bool, out []any) {
	for _, item := range out {
		nodeSet, ok := item.(map[string]any)
		if !ok {
			continue
		}
		types := toAnySlice(nodeSet["types"])
		for _, typ := range types {
			t := strings.TrimSpace(fmt.Sprint(typ))
			if t != "" {
				installed[t] = true
			}
		}
	}
}

// ReferencedNodeIDs returns all string values in nodes except each node's own id.
func ReferencedNodeIDs(nodes []map[string]any) map[string]struct{} {
	refs := make(map[string]struct{})
	for _, node := range nodes {
		collectNodeReferences(refs, node)
	}
	return refs
}

func isMissingOrphanGlobalNode(node map[string]any, installedTypes map[string]bool, allFlows []map[string]any) bool {
	if !isMissingGlobalNode(node, installedTypes) {
		return false
	}
	id := strings.TrimSpace(fmt.Sprint(node["id"]))
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
	id := strings.TrimSpace(fmt.Sprint(node["id"]))
	if id == "" {
		return false
	}
	typ := strings.TrimSpace(fmt.Sprint(node["type"]))
	if typ == "" || isBuiltInFlowType(typ) || installedTypes[typ] {
		return false
	}
	return true
}

func isMissingConfigNode(node map[string]any, installedTypes map[string]bool) bool {
	if node == nil {
		return false
	}
	id := strings.TrimSpace(fmt.Sprint(node["id"]))
	if id == "" {
		return false
	}
	typ := strings.TrimSpace(fmt.Sprint(node["type"]))
	if typ == "" || isBuiltInFlowType(typ) || installedTypes[typ] {
		return false
	}
	_, hasWires := node["wires"]
	if hasWires {
		return false
	}
	if _, hasZ := node["z"]; !hasZ {
		return true
	}
	_, hasUsers := node["_users"]
	return hasUsers
}

func isMissingNodeType(node map[string]any, installedTypes map[string]bool) bool {
	if node == nil {
		return false
	}
	id := strings.TrimSpace(fmt.Sprint(node["id"]))
	if id == "" {
		return false
	}
	typ := strings.TrimSpace(fmt.Sprint(node["type"]))
	if typ == "" || isBuiltInFlowType(typ) || installedTypes[typ] {
		return false
	}
	return true
}

func isBuiltInFlowType(typ string) bool {
	typ = strings.TrimSpace(typ)
	return typ == "tab" || typ == "group" || typ == "subflow" || strings.HasPrefix(typ, "subflow:")
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

func fetchAllFlowNodes(ctx context.Context, client *noderedclient.Client) []map[string]any {
	if client == nil {
		return nil
	}
	var out map[string]any
	code, body, errs := client.GetVersionRevV2(ctx, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		logx.WithContext(ctx).Errorf("fetch all node-red flows failed: code=%d err=%v body=%s", code, errs, string(body))
		return nil
	}
	return toMapSlice(out["flows"])
}

func fetchFlowTabs(ctx context.Context, client *noderedclient.Client) []map[string]any {
	all := fetchAllFlowNodes(ctx, client)
	tabs := make([]map[string]any, 0)
	for _, node := range all {
		if strings.TrimSpace(fmt.Sprint(node["type"])) == "tab" {
			tabs = append(tabs, node)
		}
	}
	return tabs
}

func collectNodeReferences(refs map[string]struct{}, node map[string]any) {
	if node == nil {
		return
	}
	collectReferences(refs, node, "id")
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

func fetchGlobalNodes(ctx context.Context, client *noderedclient.Client) []map[string]any {
	if client == nil {
		return nil
	}
	var out map[string]any
	code, body, errs := client.DoJSON(ctx, "GET", "/flow/global", nil, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		logx.WithContext(ctx).Errorf("fetch global flow failed: code=%d err=%v body=%s", code, errs, string(body))
		return nil
	}
	cfgs := toMapSlice(out["configs"])
	if len(cfgs) == 0 {
		cfgs = toMapSlice(out["nodes"])
	}
	return cfgs
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
