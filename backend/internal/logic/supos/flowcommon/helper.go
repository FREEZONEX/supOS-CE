package flowcommon

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/repo/relationDB"
	"backend/internal/svc"

	noderedclient "backend/share/clients/nodered"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/i18ns"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	FlowStatusDraft   = "DRAFT"
	FlowStatusPending = "PENDING"
	FlowStatusRunning = "RUNNING"
)

// FlowRepo declares the repository behaviour shared by Node-RED source/event flows.
type FlowRepo[T relationDB.FlowEntity] interface {
	FindOne(ctx context.Context, id int64) (T, error)
	Insert(ctx context.Context, data T) error
	Update(ctx context.Context, data T) error
	ReplaceModels(ctx context.Context, parentID int64, aliases []string) error
}

// FlowCopyInput defines the required fields when copying a flow.
type FlowCopyInput struct {
	FlowName    string
	Description string
	Template    string
}

// CopyFlow clones the given source flow and returns the brand new record.
func CopyFlow[T relationDB.FlowEntity](
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	repo FlowRepo[T],
	factory func() T,
	sourceID int64,
	input FlowCopyInput,
	client *noderedclient.Client,
) (T, error) {
	src, err := repo.FindOne(ctx, sourceID)
	if err != nil {
		return zero[T](), err
	}
	if isNil(src) {
		return zero[T](), errors.NotFind.WithMsg(i18ns.LocalizeMsg("nodered.flow.not.exist"))
	}

	dst := factory()
	dst.SetID(svcCtx.SnowFlake.GetSnowflakeId())
	dst.SetFlowName(strings.TrimSpace(input.FlowName))
	dst.SetDescription(strings.TrimSpace(input.Description))
	dst.SetTemplate(strings.TrimSpace(input.Template))
	// dst.SetFlowType(src.GetFlowType())
	dst.SetFlowStatus(FlowStatusDraft)

	sourceJSON, _, err := ResolveNodesJSON(ctx, client, "", src)
	if err != nil {
		return zero[T](), err
	}
	if strings.TrimSpace(sourceJSON) != "" {
		newJSON, err := regenerateNodeIDs(sourceJSON)
		if err != nil {
			return zero[T](), err
		}
		dst.SetFlowData(newJSON)
	}

	if err := repo.Insert(ctx, dst); err != nil {
		return zero[T](), err
	}
	return dst, nil
}

// DeployFlow pushes the flow definition to Node-RED and persists the latest state.
func DeployFlow[T relationDB.FlowEntity](
	ctx context.Context,
	repo FlowRepo[T],
	entityID int64,
	overrideJSON string,
	client *noderedclient.Client,
	aliasExtractor func([]map[string]any) []string,
) (string, error) {
	if client == nil {
		return "", errors.System.WithMsg(i18ns.LocalizeMsg("nodered.flow.not.exist"))
	}

	rec, err := repo.FindOne(ctx, entityID)
	if err != nil {
		return "", err
	}
	if isNil(rec) {
		return "", errors.NotFind.WithMsg(i18ns.LocalizeMsg("nodered.flow.not.exist"))
	}

	resolvedJSON, _, err := ResolveNodesJSON(ctx, client, overrideJSON, rec)
	if err != nil {
		return "", err
	}
	resolvedJSON = strings.TrimSpace(resolvedJSON)
	if resolvedJSON == "" {
		return "", errors.Parameter.WithMsg(i18ns.LocalizeMsg("nodered.flowId.empty"))
	}

	var rawNodes []map[string]any
	if err := json.Unmarshal([]byte(resolvedJSON), &rawNodes); err != nil {
		return "", errors.Parameter.WithMsg(i18ns.LocalizeMsg("nodered.invalid.parameter"))
	}

	flowNodes, globalNodes := splitGlobalNodes(rawNodes)

	flowID := strings.TrimSpace(rec.GetFlowID())
	// create flow if absent
	if flowID == "" {
		req := map[string]any{
			"id":       "",
			"nodes":    []any{},
			"disabled": false,
			"label":    rec.GetFlowName(),
			"info":     rec.GetDescription(),
		}
		var out map[string]any
		code, body, errs := client.DoJSON(ctx, "POST", "/flow", req, &out)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("create flow failed: code=%d err=%v body=%s", code, errs, string(body))
			return "", errors.System.WithMsg(i18ns.LocalizeMsg("error.sys.systemError")).AddDetailf("node-red create flow failed: code=%d err=%v body=%s", code, errs, string(body))
		}
		if id, ok := out["id"].(string); ok && strings.TrimSpace(id) != "" {
			flowID = id
		} else {
			return "", errors.System.WithMsg(i18ns.LocalizeMsg("error.sys.systemError")).AddDetail("node-red create flow returned empty id")
		}
	}

	setZ(flowNodes, flowID)

	flowBody := map[string]any{
		"id":       flowID,
		"nodes":    toInterfaceSlice(flowNodes),
		"disabled": false,
		"label":    rec.GetFlowName(),
		"info":     rec.GetDescription(),
	}
	var upd map[string]any
	code, body, errs := client.DoJSON(ctx, "PUT", "/flow/"+flowID, flowBody, &upd)
	if len(errs) > 0 || (code != 200 && code != 204) {
		logx.WithContext(ctx).Errorf("update flow failed: code=%d err=%v body=%s", code, errs, string(body))
		return "", errors.System.WithMsg(i18ns.LocalizeMsg("error.sys.systemError")).AddDetailf("node-red update flow failed: code=%d err=%v body=%s", code, errs, string(body))
	}

	if len(globalNodes) > 0 {
		globalBody := map[string]any{
			"id":      "global",
			"configs": toInterfaceSlice(globalNodes),
		}
		var gout map[string]any
		code, body, errs = client.DoJSON(ctx, "PUT", "/flow/global", globalBody, &gout)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("update global flow failed: code=%d err=%v body=%s", code, errs, string(body))
			return "", errors.System.WithMsg(i18ns.LocalizeMsg("error.sys.systemError")).AddDetailf("node-red update global failed: code=%d err=%v body=%s", code, errs, string(body))
		}
	}

	merged := append(flowNodes, globalNodes...)
	newFlowData, err := json.Marshal(merged)
	if err != nil {
		return "", errors.System.WithMsg(err.Error())
	}

	rec.SetFlowID(flowID)
	rec.SetFlowData(string(newFlowData))
	rec.SetFlowStatus(FlowStatusRunning)

	if err := repo.Update(ctx, rec); err != nil {
		return "", err
	}

	if aliasExtractor != nil {
		aliases := aliasExtractor(rawNodes)
		if err := repo.ReplaceModels(ctx, rec.GetID(), aliases); err != nil {
			return "", err
		}
	} else {
		_ = repo.ReplaceModels(ctx, rec.GetID(), nil)
	}

	return flowID, nil
}

// ResolveNodesJSON resolves the JSON array describing the Node-RED nodes with the precedence:
// override > draft(FlowData) > Node-RED runtime.
func ResolveNodesJSON(ctx context.Context, client *noderedclient.Client, override string, entity relationDB.FlowEntity) (string, string, error) {
	if strings.TrimSpace(override) != "" {
		return override, "override", nil
	}
	if data := strings.TrimSpace(entity.GetFlowData()); data != "" {
		return data, "draft", nil
	}
	if client != nil && strings.TrimSpace(entity.GetFlowID()) != "" {
		var out map[string]any
		code, body, errs := client.GetFlowNodesV1(ctx, entity.GetFlowID(), &out)
		if len(errs) > 0 || (code != 200 && code != 204) {
			logx.WithContext(ctx).Errorf("fetch nodes from node-red failed: code=%d err=%v body=%s", code, errs, string(body))
			return "", "", errors.System.WithMsg(i18ns.LocalizeMsg("nodered.flow.not.exist"))
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
		return "", errors.Parameter.WithMsg(i18ns.LocalizeMsg("nodered.invalid.parameter"))
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

func zero[T any]() T {
	var v T
	return v
}

func isNil[T relationDB.FlowEntity](v T) bool {
	return any(v) == nil
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
