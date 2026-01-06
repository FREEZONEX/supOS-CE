package flowcommon

import (
	"context"
	"encoding/json"
	"fmt"
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

	flowNodes, globalNodes := splitGlobalNodes(rawNodes)

	// mqtt的全局节点先部署
	if len(globalNodes) > 0 {
		mergedGlobal := mergeGlobalNodes(ctx, client, globalNodes)
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

	syncSupmodelMappings(ctx, client, flowNodes)

	if len(flowNodes) > 0 {
		aliases := aliasExtractor(flowNodes)
		if err := repo.ReplaceModels(ctx, entityID, aliases); err != nil {
			return "", err
		}
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

func syncSupmodelMappings(ctx context.Context, client *noderedclient.Client, nodes []map[string]any) {
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if strings.TrimSpace(fmt.Sprint(node["type"])) != "supmodel" {
			continue
		}
		nodeID := strings.TrimSpace(fmt.Sprint(node["id"]))
		if nodeID == "" {
			continue
		}
		mapping := toAnySlice(node["mapping"])
		if len(mapping) == 0 {
			continue
		}
		req := map[string]any{
			"nodeId":  nodeID,
			"mapping": mapping,
		}
		var out map[string]any
		code, body, errs := client.DoJSON(ctx, "POST", "/nodered-api/upload/tags", req, &out)
		if len(errs) > 0 || (code != 200 && code != 201) {
			logx.WithContext(ctx).Errorf("syncSupmodelMappings sync supmodel mapping failed: id=%s code=%d err=%v body=%s", nodeID, code, errs, string(body))
		}
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

// ExtractAliases parses possible UNS aliases from Node-RED node definitions.
func ExtractAliases(nodes []map[string]any) []string {
	aliasSet := make(map[string]struct{})
	for _, node := range nodes {
		if node == nil {
			continue
		}
		for _, key := range []string{"selectedAlias", "selectedModelAlias", "alias", "unsAlias"} {
			if val, ok := node[key]; ok {
				if alias := strings.TrimSpace(fmt.Sprint(val)); alias != "" {
					aliasSet[alias] = struct{}{}
				}
			}
		}
		if val, ok := node["selectedModel"]; ok {
			if alias := strings.TrimSpace(fmt.Sprint(val)); alias != "" {
				aliasSet[alias] = struct{}{}
			}
		}
	}
	aliases := make([]string, 0, len(aliasSet))
	for alias := range aliasSet {
		aliases = append(aliases, alias)
	}
	return aliases
}

func mergeGlobalNodes(ctx context.Context, client *noderedclient.Client, incoming []map[string]any) []map[string]any {
	existing := fetchGlobalNodes(ctx, client)
	if len(existing) == 0 {
		return incoming
	}

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
