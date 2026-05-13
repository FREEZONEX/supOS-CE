package flowcommon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/common"
	"backend/internal/repo/relationDB"
	noderedclient "backend/share/clients/nodered"
)

type fakeFlowRepo struct {
	source   *relationDB.NoderedFlow
	inserted *relationDB.NoderedFlow
}

func (f *fakeFlowRepo) FindOne(_ context.Context, id int64) (*relationDB.NoderedFlow, error) {
	if f.source != nil && f.source.ID == id {
		return f.source, nil
	}
	return nil, nil
}

func (f *fakeFlowRepo) Insert(_ context.Context, data *relationDB.NoderedFlow) error {
	f.inserted = data
	return nil
}

func (f *fakeFlowRepo) Update(_ context.Context, _ *relationDB.NoderedFlow) error {
	return nil
}

func (f *fakeFlowRepo) ReplaceModels(_ context.Context, _ int64, _ []string) error {
	return nil
}

func TestCopyFlowSetsCreatorFromInput(t *testing.T) {
	common.InitSnowflake(1)
	repo := &fakeFlowRepo{
		source: &relationDB.NoderedFlow{
			ID:       1001,
			FlowData: `[{"id":"node-1","z":"tab-1","type":"debug"}]`,
		},
	}

	record, err := CopyFlow(context.Background(), nil, repo, 1001, FlowCopyInput{
		FlowName: "copied-flow",
		Template: "node-red",
		Creator:  "tier0",
	}, nil)
	if err != nil {
		t.Fatalf("CopyFlow returned error: %v", err)
	}
	if record == nil {
		t.Fatalf("CopyFlow returned nil record")
	}
	if record.Creator != "tier0" {
		t.Fatalf("expected creator to be tier0, got %q", record.Creator)
	}
	if repo.inserted == nil {
		t.Fatalf("expected inserted record to be captured")
	}
	if repo.inserted.Creator != "tier0" {
		t.Fatalf("expected inserted creator to be tier0, got %q", repo.inserted.Creator)
	}
}

func TestReferencedNodeIDsSkipsOwnID(t *testing.T) {
	refs := ReferencedNodeIDs([]map[string]any{
		{
			"id":     "node-1",
			"type":   "mqtt out",
			"broker": "broker-1",
			"wires":  []any{[]any{"debug-1"}},
		},
	})

	if _, ok := refs["node-1"]; ok {
		t.Fatalf("own id should not be considered a reference")
	}
	if _, ok := refs["broker-1"]; !ok {
		t.Fatalf("expected broker id to be collected")
	}
	if _, ok := refs["debug-1"]; !ok {
		t.Fatalf("expected nested wire id to be collected")
	}
}

func TestParseInstalledNodeTypesFromNodeEditorHTML(t *testing.T) {
	types := parseInstalledNodeTypes([]byte(`
		<script type="text/javascript">
			RED.nodes.registerType('inject',{});
			RED.nodes.registerType("mqtt-broker",{});
		</script>
	`))

	if !types["inject"] {
		t.Fatalf("expected inject type to be parsed")
	}
	if !types["mqtt-broker"] {
		t.Fatalf("expected mqtt-broker type to be parsed")
	}
	if types["missing-type"] {
		t.Fatalf("unexpected missing type")
	}
}

func TestMissingNodeTypesReportsNormalAndConfigNodes(t *testing.T) {
	installed := map[string]bool{
		"inject":      true,
		"mqtt-broker": true,
	}
	missing := MissingNodeTypes([]map[string]any{
		{
			"id":   "bad-config",
			"type": "missing-config-type",
		},
		{
			"id":   "bad-normal-node",
			"type": "missing-node-type",
			"z":    "flow-1",
		},
		{
			"id":      "bad-normal-terminal",
			"type":    "missing-terminal-node",
			"z":       "flow-1",
			"wires":   []any{},
			"outputs": 0,
		},
		{
			"id":     "bad-tab-config",
			"type":   "missing-tab-config-type",
			"z":      "flow-1",
			"_users": []any{"node-1"},
		},
		{
			"id":   "broker-1",
			"type": "mqtt-broker",
		},
		{
			"id":   "group-1",
			"type": "group",
			"z":    "flow-1",
		},
		{
			"id":   "subflow-def",
			"type": "subflow",
		},
		{
			"id":   "subflow-inst",
			"type": "subflow:subflow-def",
			"z":    "flow-1",
		},
	}, installed)

	expected := []string{"missing-config-type", "missing-node-type", "missing-tab-config-type", "missing-terminal-node"}
	if len(missing) != len(expected) {
		t.Fatalf("expected missing types %v, got %#v", expected, missing)
	}
	for i, typ := range expected {
		if missing[i] != typ {
			t.Fatalf("expected missing types %v, got %#v", expected, missing)
		}
	}
}

func TestRemoveWireReferencesDropsRemovedTargets(t *testing.T) {
	nodes := []map[string]any{
		{
			"id":    "node-1",
			"wires": []any{[]any{"keep-1", "drop-1"}, []any{"drop-2"}},
		},
	}
	removeWireReferences(nodes, map[string]struct{}{
		"drop-1": {},
		"drop-2": {},
	})
	wires := nodes[0]["wires"].([]any)
	first := wires[0].([]any)
	second := wires[1].([]any)
	if len(first) != 1 || first[0] != "keep-1" || len(second) != 0 {
		t.Fatalf("unexpected wires after cleanup: %#v", wires)
	}
}

func TestIsMissingOrphanGlobalNode(t *testing.T) {
	installed := map[string]bool{
		"mqtt-broker": true,
	}
	orphan := map[string]any{
		"id":   "bad-config",
		"type": "missing-type",
	}
	referenced := map[string]any{
		"id":     "mqtt-out-1",
		"type":   "mqtt out",
		"broker": "bad-config",
		"z":      "flow-1",
	}
	if !isMissingOrphanGlobalNode(orphan, installed, []map[string]any{orphan}) {
		t.Fatalf("expected missing unreferenced global node to be orphan")
	}
	if isMissingOrphanGlobalNode(orphan, installed, []map[string]any{orphan, referenced}) {
		t.Fatalf("referenced missing global node should not be treated as orphan")
	}
	if isMissingOrphanGlobalNode(map[string]any{"id": "broker-1", "type": "mqtt-broker"}, installed, nil) {
		t.Fatalf("installed global node should not be treated as orphan")
	}
}

func TestListMissingRuntimeNodesFindsGlobalAndHiddenFlowNodes(t *testing.T) {
	client, _ := newNodeRedRuntimeTestServer(t)

	nodes, err := ListMissingRuntimeNodes(context.Background(), client)
	if err != nil {
		t.Fatalf("ListMissingRuntimeNodes returned error: %v", err)
	}
	got := make(map[string]RuntimeMissingNode)
	for _, node := range nodes {
		got[node.ID] = node
	}
	if got["bad-global"].Scope != MissingNodeScopeGlobalConfig || got["bad-global"].FlowID != "global" {
		t.Fatalf("expected bad global config to be listed, got %#v", got["bad-global"])
	}
	if got["bad-node"].Scope != MissingNodeScopeFlowNode || got["bad-node"].FlowID != "tab-2" {
		t.Fatalf("expected hidden bad flow node to be listed, got %#v", got["bad-node"])
	}
	if got["bad-config"].Scope != MissingNodeScopeFlowConfig || got["bad-config"].FlowID != "tab-2" {
		t.Fatalf("expected hidden bad flow config to be listed, got %#v", got["bad-config"])
	}
}

func TestDeleteMissingRuntimeNodeRemovesFlowNodeAndWires(t *testing.T) {
	client, state := newNodeRedRuntimeTestServer(t)

	deleted, err := DeleteMissingRuntimeNode(context.Background(), client, MissingNodeDeleteTarget{
		ID:     "bad-node",
		FlowID: "tab-2",
		Scope:  MissingNodeScopeFlowNode,
	})
	if err != nil {
		t.Fatalf("DeleteMissingRuntimeNode returned error: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected one node deleted, got %d", deleted)
	}
	for _, node := range state.flows["tab-2"].nodes {
		if node["id"] == "bad-node" {
			t.Fatalf("bad node was not removed: %#v", state.flows["tab-2"].nodes)
		}
		if node["id"] == "inject-2" {
			wires := node["wires"].([]any)
			first := wires[0].([]any)
			if len(first) != 1 || first[0] != "debug-2" {
				t.Fatalf("expected wire to bad-node to be removed, got %#v", wires)
			}
		}
	}
}

func TestDeleteMissingRuntimeNodeRemovesGlobalConfig(t *testing.T) {
	client, state := newNodeRedRuntimeTestServer(t)

	deleted, err := DeleteMissingRuntimeNode(context.Background(), client, MissingNodeDeleteTarget{
		ID:     "bad-global",
		FlowID: "global",
		Scope:  MissingNodeScopeGlobalConfig,
	})
	if err != nil {
		t.Fatalf("DeleteMissingRuntimeNode returned error: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected one global config deleted, got %d", deleted)
	}
	for _, node := range state.global {
		if node["id"] == "bad-global" {
			t.Fatalf("bad global config was not removed: %#v", state.global)
		}
	}
}

type nodeRedRuntimeTestState struct {
	global []map[string]any
	flows  map[string]*nodeRedRuntimeTestFlow
}

type nodeRedRuntimeTestFlow struct {
	label   string
	nodes   []map[string]any
	configs []map[string]any
}

func newNodeRedRuntimeTestServer(t *testing.T) (*noderedclient.Client, *nodeRedRuntimeTestState) {
	t.Helper()
	state := &nodeRedRuntimeTestState{
		global: []map[string]any{
			{"id": "broker-1", "type": "mqtt-broker", "name": "emqx:1883"},
			{"id": "bad-global", "type": "missing-global-type", "name": "bad global"},
		},
		flows: map[string]*nodeRedRuntimeTestFlow{
			"tab-1": {
				label: "Visible",
				nodes: []map[string]any{
					{"id": "inject-1", "type": "inject", "z": "tab-1", "wires": []any{[]any{"debug-1"}}},
					{"id": "debug-1", "type": "debug", "z": "tab-1", "wires": []any{}},
				},
			},
			"tab-2": {
				label: "Hidden",
				nodes: []map[string]any{
					{"id": "inject-2", "type": "inject", "z": "tab-2", "wires": []any{[]any{"bad-node", "debug-2"}}},
					{"id": "bad-node", "type": "missing-node-type", "z": "tab-2", "wires": []any{[]any{"debug-2"}}},
					{"id": "debug-2", "type": "debug", "z": "tab-2", "wires": []any{}},
				},
				configs: []map[string]any{
					{"id": "bad-config", "type": "missing-config-type", "z": "tab-2", "_users": []any{"inject-2"}},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/nodes":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`RED.nodes.registerType('inject',{});RED.nodes.registerType('debug',{});RED.nodes.registerType('mqtt-broker',{});`))
		case r.Method == http.MethodGet && r.URL.Path == "/flow/global":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "global", "configs": state.global})
		case r.Method == http.MethodPut && r.URL.Path == "/flow/global":
			var req map[string]any
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("decode global put: %v", err)
			}
			state.global = toMapSlice(req["configs"])
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "global"})
		case r.Method == http.MethodGet && r.URL.Path == "/flows":
			flows := []map[string]any{
				{"id": "tab-1", "type": "tab", "label": "Visible"},
				{"id": "tab-2", "type": "tab", "label": "Hidden"},
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"flows": flows, "rev": "1"})
		case r.Method == http.MethodGet && r.URL.Path == "/flow/tab-1":
			writeRuntimeTestFlow(w, "tab-1", state.flows["tab-1"])
		case r.Method == http.MethodGet && r.URL.Path == "/flow/tab-2":
			writeRuntimeTestFlow(w, "tab-2", state.flows["tab-2"])
		case r.Method == http.MethodPut && r.URL.Path == "/flow/tab-1":
			readRuntimeTestFlow(t, r, state.flows["tab-1"])
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "tab-1"})
		case r.Method == http.MethodPut && r.URL.Path == "/flow/tab-2":
			readRuntimeTestFlow(t, r, state.flows["tab-2"])
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "tab-2"})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	return noderedclient.NewClient(server.URL, "", ""), state
}

func writeRuntimeTestFlow(w http.ResponseWriter, id string, flow *nodeRedRuntimeTestFlow) {
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":       id,
		"label":    flow.label,
		"disabled": false,
		"info":     "",
		"nodes":    flow.nodes,
		"configs":  flow.configs,
	})
}

func readRuntimeTestFlow(t *testing.T, r *http.Request, flow *nodeRedRuntimeTestFlow) {
	t.Helper()
	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.Fatalf("decode flow put: %v", err)
	}
	flow.nodes = toMapSlice(req["nodes"])
	flow.configs = toMapSlice(req["configs"])
}
