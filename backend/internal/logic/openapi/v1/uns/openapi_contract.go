package uns

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/dataingest"
	domainuns "backend/internal/domain/uns"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/repo"
	"backend/internal/svc"
	"backend/internal/types"
)

const (
	openapiQualityGood        = "Good"
	openapiQualityUncertain   = "Uncertain"
	openapiQualityBad         = "Bad"
	openapiQualityGoodNoData  = "GoodNoData"
	openapiHistoryDefaultSize = int64(1000)
	openapiHistoryMaxSize     = int64(2000)
)

func openapiUnsRead(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ReadReq) (*types.Envelope, error) {
	if req == nil || len(req.Topics) == 0 {
		return respx.Envelope(map[string]any{
			"success": false,
			"results": []map[string]any{{"success": false, "error": openapiItemError(400, "topics empty")}},
		}), nil
	}
	topics, err := expandOpenapiTopics(ctx, req.Topics)
	if err != nil {
		return nil, logicx.Error(err)
	}
	nodeByTopic := make(map[string]repo.UnsNode, len(topics))
	metadataNodes := make([]repo.UnsNode, 0, len(topics))
	for _, topic := range topics {
		node, ok, err := findOpenapiNode(ctx, topic, false)
		if err != nil {
			return nil, logicx.Error(err)
		}
		if !ok {
			continue
		}
		nodeByTopic[topic] = node
		if req.IncludeMetadata {
			metadataNodes = append(metadataNodes, node)
		}
	}
	bindFlowIDs, err := openapiBindFlowIDMap(ctx, req.IncludeMetadata, metadataNodes)
	if err != nil {
		return nil, logicx.Error(err)
	}
	results := make([]map[string]any, 0, len(topics))
	success := true
	for _, topic := range topics {
		row := map[string]any{"success": false, "topic": topic}
		node, ok := nodeByTopic[topic]
		if !ok {
			row["error"] = openapiItemError(404, "topic not found")
			success = false
			results = append(results, row)
			continue
		}
		if node.Type != 2 {
			row["error"] = openapiItemError(400, "topic is not a file node")
			success = false
			results = append(results, row)
			continue
		}
		row["success"] = true
		row["result"] = openapiLatestVQT(ctx, svcCtx, node)
		if req.IncludeMetadata {
			row["metadata"] = openapiNodeInfo(ctx, svcCtx, node, true, req.IncludeLeafValue, bindFlowIDs)
		}
		results = append(results, row)
	}
	return respx.Envelope(map[string]any{"success": success, "results": results}), nil
}

func openapiUnsWrite(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WriteReq) (*types.Envelope, error) {
	if req == nil || len(req.Writes) == 0 {
		return nil, logicx.Error(fmt.Errorf("writes is required"))
	}
	if req.Qos < 0 || req.Qos > 2 {
		return nil, logicx.Error(fmt.Errorf("qos must be 0, 1 or 2"))
	}
	results := make([]map[string]any, 0, len(req.Writes))
	success := true
	for _, item := range req.Writes {
		topic := normalizeOpenapiTopic(item.Topic)
		row := map[string]any{"success": false, "topic": topic}
		switch {
		case topic == "":
			row["error"] = openapiItemError(400, "topic empty")
		case hasOpenapiWildcard(topic):
			row["error"] = openapiItemError(400, "write topic does not support wildcard")
		default:
			node, ok, err := findOpenapiNode(ctx, topic, false)
			if err != nil {
				return nil, logicx.Error(err)
			}
			if !ok {
				row["error"] = openapiItemError(404, "topic not found")
			} else if node.Type != 2 {
				row["error"] = openapiItemError(400, "topic is not a file node")
			} else if svcCtx == nil || svcCtx.App == nil || svcCtx.App.DataIngest == nil {
				row["error"] = openapiItemError(500, "data ingest is not initialized")
			} else if err := svcCtx.App.DataIngest.WriteOpenAPIValue(
				ctx,
				node,
				item.Value,
				item.TimeStamp,
				byte(req.Qos),
				req.Retain,
			); err != nil {
				row["error"] = openapiItemError(500, err.Error())
			} else {
				row["success"] = true
			}
		}
		if row["success"] != true {
			success = false
		}
		results = append(results, row)
	}
	return respx.Envelope(map[string]any{"success": success, "results": results}), nil
}

func openapiUnsBrowse(ctx context.Context, svcCtx *svc.ServiceContext, req *types.BrowseReq) (*types.Envelope, error) {
	if req == nil {
		req = &types.BrowseReq{}
	}
	nodes, root, err := openapiBrowseNodes(ctx, req.Path)
	if err != nil {
		return nil, logicx.Error(err)
	}
	bindFlowIDs, err := openapiBindFlowIDMap(ctx, req.IncludeMetadata, nodes)
	if err != nil {
		return nil, logicx.Error(err)
	}
	tree := buildOpenapiNodeTree(ctx, svcCtx, nodes, root, req.MaxDepth, req.IncludeMetadata, req.IncludeLeafValue, bindFlowIDs)
	return respx.Envelope(map[string]any{"tree": tree}), nil
}

func openapiUnsSearch(ctx context.Context, svcCtx *svc.ServiceContext, req *types.SearchReq) (*types.Envelope, error) {
	if req == nil {
		req = &types.SearchReq{}
	}
	nodes, err := repo.NewUnsRepo(ctx).ListUnsNodes(ctx, repo.UnsNodeFilter{Keyword: strings.TrimSpace(req.Keyword)})
	if err != nil {
		return nil, logicx.Error(err)
	}
	prefix := normalizeOpenapiTopic(req.PathPrefix)
	topicType, topicTypeSet, err := openapiTopicTypeValue(req.TopicType)
	if err != nil {
		return nil, logicx.Error(err)
	}
	filtered := make([]repo.UnsNode, 0, len(nodes))
	for _, node := range nodes {
		if prefix != "" && node.Namespace != prefix && !strings.HasPrefix(node.Namespace, prefix+"/") {
			continue
		}
		if topicTypeSet {
			if node.TopicType != topicType {
				continue
			}
		}
		filtered = append(filtered, node)
	}
	sort.SliceStable(filtered, func(i, j int) bool { return filtered[i].Namespace < filtered[j].Namespace })
	page, size := normalizeOpenapiPage(req.Page, req.Size, 20)
	total := int64(len(filtered))
	filtered = pageOpenapiNodes(filtered, page, size)
	bindFlowIDs, err := openapiBindFlowIDMap(ctx, req.IncludeMetadata, filtered)
	if err != nil {
		return nil, logicx.Error(err)
	}
	objects := make([]map[string]any, 0, len(filtered))
	for _, node := range filtered {
		objects = append(objects, openapiNodeInfo(ctx, svcCtx, node, req.IncludeMetadata, req.IncludeLeafValue, bindFlowIDs))
	}
	return respx.Envelope(map[string]any{
		"objects": objects,
		"page":    page,
		"size":    size,
		"total":   total,
	}), nil
}

func openapiUnsHistory(ctx context.Context, svcCtx *svc.ServiceContext, req *types.HistoryReq) (*types.Envelope, error) {
	if req == nil || len(req.Topics) == 0 {
		return nil, logicx.Error(fmt.Errorf("topics is required"))
	}
	start, err := parseOpenapiHistoryTime("start_time", req.StartTime)
	if err != nil {
		return nil, logicx.Error(err)
	}
	end, err := parseOpenapiHistoryTime("end_time", req.EndTime)
	if err != nil {
		return nil, logicx.Error(err)
	}
	if end < start {
		return nil, logicx.Error(fmt.Errorf("end_time must be greater than or equal to start_time"))
	}
	agg, aggErr := parseOpenapiHistoryAggregation(req.Aggregation)
	autoMode := req.Page <= 0 && req.Size <= 0
	page, size := normalizeOpenapiPage(req.Page, req.Size, openapiHistoryDefaultSize)
	if autoMode {
		page = 1
		size = openapiHistoryDefaultSize
	}
	results := make([]map[string]any, 0, len(req.Topics))
	success := true
	total := int64(0)
	for _, rawTopic := range req.Topics {
		topic := normalizeOpenapiTopic(rawTopic)
		row := map[string]any{"success": false, "topic": topic}
		switch {
		case topic == "":
			row["error"] = openapiItemError(400, "topic empty")
		case hasOpenapiWildcard(topic):
			row["error"] = openapiItemError(400, "history topic does not support wildcard")
		default:
			node, ok, err := findOpenapiNode(ctx, topic, false)
			if err != nil {
				return nil, logicx.Error(err)
			}
			if !ok {
				row["error"] = openapiItemError(404, "topic not found")
			} else if node.Type != 2 {
				row["error"] = openapiItemError(400, "topic is not a file node")
			} else if node.EnableHistory != 1 {
				row["error"] = openapiItemError(400, "history is not enabled for topic")
			} else if aggErr != nil {
				row["error"] = openapiItemError(400, aggErr.Error())
			} else if svcCtx == nil || svcCtx.App == nil || svcCtx.App.DataIngest == nil {
				row["error"] = openapiItemError(500, "data ingest is not initialized")
			} else {
				values, itemTotal, meta, err := queryOpenapiHistory(ctx, svcCtx, node, start, end, page, size, autoMode, agg)
				if err != nil {
					row["error"] = openapiItemError(400, err.Error())
				} else {
					row["success"] = true
					row["result"] = map[string]any{"values": values, "meta": meta}
					total += itemTotal
				}
			}
		}
		if row["success"] != true {
			success = false
		}
		results = append(results, row)
	}
	return respx.Envelope(map[string]any{
		"success": success,
		"results": results,
		"page":    page,
		"size":    size,
		"total":   total,
	}), nil
}

func openapiUnsCreate(ctx context.Context, svcCtx *svc.ServiceContext, req *types.NodeCreateReq) (*types.Envelope, error) {
	if req == nil || len(req.Namespace) == 0 {
		return nil, logicx.Error(fmt.Errorf("namespace is required"))
	}
	roots := make([]domainuns.CreateTreeNode, 0, len(req.Namespace))
	userID := logicx.UserID(ctx)
	for _, node := range req.Namespace {
		roots = append(roots, openapiCreateTreeNode(node, "", userID))
	}
	created, err := svcCtx.App.UNS.CreateTreeBatch(ctx, roots)
	if err != nil {
		return nil, logicx.Error(err)
	}
	results := make([]map[string]any, 0, len(created))
	success := true
	for _, item := range created {
		row := map[string]any{"success": item.Err == nil, "topic": item.Name}
		if item.Err != nil {
			row["error"] = openapiItemError(400, item.Err.Error())
			success = false
		} else {
			row["topic"] = item.Namespace
		}
		results = append(results, row)
	}
	return respx.Envelope(map[string]any{"success": success, "results": results}), nil
}

func openapiUnsUpdate(ctx context.Context, svcCtx *svc.ServiceContext, req *types.NodeUpdateReq) (*types.Envelope, error) {
	if req == nil || strings.TrimSpace(req.Path) == "" {
		return nil, logicx.Error(fmt.Errorf("path is required"))
	}
	node, ok, err := findOpenapiNode(ctx, req.Path, false)
	if err != nil {
		return nil, logicx.Error(err)
	}
	if !ok {
		return nil, logicx.Error(domainuns.ErrNotFound)
	}
	name := firstNonBlank(req.Name, node.Name)
	displayName := firstNonBlank(req.DisplayName, node.DisplayName)
	description := firstNonBlank(req.Description, node.Description)
	schema := string(node.Schema)
	if len(req.Fields) > 0 {
		schema = openapiSchemaFieldsJSON(req.Fields)
	}
	extend := string(node.ExtendProperties)
	if len(req.ExtendProperties) > 0 {
		extend = openapiStringMapJSON(req.ExtendProperties)
	}
	enableHistory := node.EnableHistory == 1
	if req.EnableHistory != nil {
		enableHistory = *req.EnableHistory
	}
	_, err = svcCtx.App.UNS.Update(ctx, domainuns.SaveCommand{
		ID:               node.ID,
		ParentID:         node.ParentID,
		Name:             name,
		DisplayName:      displayName,
		Description:      description,
		Alias:            firstNonBlank(req.Alias, node.Alias),
		NodeType:         openapiInternalNodeType(node.Type),
		TopicType:        openapiInternalTopicType(node.TopicType),
		Schema:           schema,
		ExtendProperties: extend,
		EnableHistory:    &enableHistory,
		UserID:           logicx.UserID(ctx),
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(map[string]any{"success": true}), nil
}

func openapiUnsDelete(ctx context.Context, svcCtx *svc.ServiceContext, req *types.NodeDeleteReq) (*types.Envelope, error) {
	if req == nil || len(req.Topics) == 0 { return nil, logicx.Error(fmt.Errorf("topics is required")) }
	for _, rawTopic := range req.Topics {
		topic := normalizeOpenapiTopic(rawTopic)
		if topic == "" || hasOpenapiWildcard(topic) { return nil, logicx.Error(fmt.Errorf("delete topics does not support wildcard")) }
		node, ok, err := findOpenapiNode(ctx, topic, false)
		if err != nil { return nil, logicx.Error(err) }
		if !ok { return nil, logicx.Error(domainuns.ErrNotFound) }
		if _, err := svcCtx.App.UNS.Delete(ctx, node.ID, logicx.UserID(ctx)); err != nil { return nil, logicx.Error(err) }
	}
	return respx.Envelope(map[string]any{"success": true}), nil
}


func openapiBrowseNodes(ctx context.Context, path string) ([]repo.UnsNode, *repo.UnsNode, error) {
	path = normalizeOpenapiTopic(path)
	unsRepo := repo.NewUnsRepo(ctx)
	if path == "" {
		nodes, err := unsRepo.ListUnsNodes(ctx, repo.UnsNodeFilter{})
		return nodes, nil, err
	}
	root, ok, err := findOpenapiNode(ctx, path, false)
	if err != nil || !ok {
		return nil, nil, err
	}
	nodes, err := unsRepo.ListUnsSubtree(ctx, root.ID)
	return nodes, &root, err
}

func buildOpenapiNodeTree(ctx context.Context, svcCtx *svc.ServiceContext, nodes []repo.UnsNode, root *repo.UnsNode, maxDepth int64, includeMetadata, includeLeafValue bool, bindFlowIDs map[int64]int64) []map[string]any {
	byParent := make(map[int64][]repo.UnsNode)
	for _, node := range nodes {
		byParent[node.ParentID] = append(byParent[node.ParentID], node)
	}
	var walk func(node repo.UnsNode, depth int64) map[string]any
	walk = func(node repo.UnsNode, depth int64) map[string]any {
		info := openapiNodeInfo(ctx, svcCtx, node, includeMetadata, includeLeafValue, bindFlowIDs)
		if maxDepth <= 0 || depth < maxDepth {
			children := byParent[node.ID]
			if len(children) > 0 {
				out := make([]map[string]any, 0, len(children))
				for _, child := range children {
					out = append(out, walk(child, depth+1))
				}
				info["children"] = out
			}
		}
		return info
	}
	if root != nil {
		return []map[string]any{walk(*root, 1)}
	}
	roots := make([]repo.UnsNode, 0)
	for _, node := range byParent[0] {
		roots = append(roots, node)
	}
	sort.SliceStable(roots, func(i, j int) bool { return roots[i].SortKey < roots[j].SortKey })
	tree := make([]map[string]any, 0, len(roots))
	for _, node := range roots {
		tree = append(tree, walk(node, 1))
	}
	return tree
}

func openapiNodeInfo(ctx context.Context, svcCtx *svc.ServiceContext, node repo.UnsNode, includeMetadata, includeLeafValue bool, bindFlowIDs ...map[int64]int64) map[string]any {
	ret := map[string]any{
		"id":        node.ID,
		"path":      strings.Trim(node.Namespace, "/"),
		"name":      node.Name,
		"type":      openapiExternalNodeType(node.Type),
		"topicType": openapiExternalTopicType(node.TopicType),
	}
	if includeMetadata {
		ret["description"] = node.Description
		ret["displayName"] = node.DisplayName
		ret["extendProperties"] = openapiExtendProperties(node.ExtendProperties)
		ret["fields"] = openapiSchemaFields(node.Schema)
		ret["alias"] = strings.Trim(node.Alias, "/")
		ret["enableHistory"] = node.EnableHistory == 1
		ret["bindFlowID"] = int64(0)
		if len(bindFlowIDs) > 0 && bindFlowIDs[0] != nil {
			ret["bindFlowID"] = bindFlowIDs[0][node.ID]
		}
	}
	if includeLeafValue && node.Type == 2 {
		ret["payload"] = openapiLatestVQT(ctx, svcCtx, node)
	}
	return ret
}

func openapiBindFlowIDMap(ctx context.Context, includeMetadata bool, nodes []repo.UnsNode) (map[int64]int64, error) {
	out := map[int64]int64{}
	if !includeMetadata || len(nodes) == 0 {
		return out, nil
	}
	ids := make([]int64, 0, len(nodes))
	seen := make(map[int64]struct{}, len(nodes))
	for _, node := range nodes {
		if node.ID <= 0 {
			continue
		}
		if _, ok := seen[node.ID]; ok {
			continue
		}
		seen[node.ID] = struct{}{}
		ids = append(ids, node.ID)
	}
	return repo.NewFlowRepo(ctx).BindFlowIDMapByNodeIDs(ctx, ids)
}

func openapiLatestVQT(ctx context.Context, svcCtx *svc.ServiceContext, node repo.UnsNode) map[string]any {
	ret := map[string]any{"quality": openapiQualityGoodNoData, "timeStamp": int64(0)}
	if svcCtx == nil || svcCtx.App == nil || svcCtx.App.DataIngest == nil || node.Type != 2 {
		return ret
	}
	payload, err := svcCtx.App.DataIngest.LastPayload(ctx, node)
	if err != nil || len(payload) == 0 {
		return ret
	}
	value := payload["payload"]
	if value == nil {
		value = payload["data"]
	}
	ret["value"] = openapiLatestValuePayload(value)
	ret["quality"] = openapiPayloadQuality(payload["quality"])
	ret["timeStamp"] = openapiLatestTimestamp(payload)
	return ret
}

func openapiCreateTreeNode(node types.NamespaceNode, inheritedTopicType string, userID int64) domainuns.CreateTreeNode {
	name := strings.TrimSpace(node.Name)
	nodeType := normalizeNamespaceNodeType(node.Type, len(node.Children) > 0)
	topicType := strings.TrimSpace(node.TopicType)
	if topicType == "" {
		topicType = inheritedTopicType
	}
	if nodeType == "folder" {
		if category := openapiTopicCategoryByName(name); category != "" {
			topicType = category
		}
	}
	enableHistory := true
	if node.EnableHistory != nil {
		enableHistory = *node.EnableHistory
	}
	result := domainuns.CreateTreeNode{Command: domainuns.SaveCommand{
		Name:             name,
		DisplayName:      node.DisplayName,
		Description:      node.Description,
		Alias:            node.Alias,
		NodeType:         nodeType,
		TopicType:        topicType,
		Schema:           openapiSchemaFieldsJSON(node.Fields),
		ExtendProperties: openapiStringMapJSON(node.ExtendProperties),
		EnableHistory:    &enableHistory,
		UserID:           userID,
	}}
	if nodeType != "folder" {
		return result
	}
	nextInherited := inheritedTopicType
	if topicType != "" {
		nextInherited = topicType
	}
	result.Children = make([]domainuns.CreateTreeNode, 0, len(node.Children))
	for _, child := range node.Children {
		result.Children = append(result.Children, openapiCreateTreeNode(child, nextInherited, userID))
	}
	return result
}

func queryOpenapiHistory(ctx context.Context, svcCtx *svc.ServiceContext, node repo.UnsNode, start, end, page, size int64, autoMode bool, agg openapiHistoryAggregation) ([]map[string]any, int64, map[string]any, error) {
	query := dataingest.HistoryQuery{
		StartMs:    start,
		EndMs:      end,
		Page:       int(page),
		Size:       int(size),
		Auto:       autoMode,
		Aggregate:  agg.enabled,
		IntervalMs: agg.intervalMs,
	}
	if agg.legacy {
		query.Fields = []dataingest.HistoryAggregationField{{Name: agg.field, Function: agg.function}}
	} else if len(agg.fields) > 0 {
		query.Fields = make([]dataingest.HistoryAggregationField, 0, len(agg.fields))
		for _, field := range agg.fields {
			query.Fields = append(query.Fields, dataingest.HistoryAggregationField{Name: field.name, Function: field.function})
		}
	}
	data, err := svcCtx.App.DataIngest.QueryHistory(ctx, node, query)
	if err != nil {
		return nil, 0, nil, err
	}
	values := openapiHistoryValuesFromData(data, agg)
	meta, _ := data["meta"].(map[string]any)
	if meta == nil {
		meta = map[string]any{}
	}
	return values, openapiAnyInt64(data["total"]), meta, nil
}

func openapiHistoryValuesFromData(data map[string]any, agg openapiHistoryAggregation) []map[string]any {
	list, _ := data["list"].([]map[string]any)
	values := make([]map[string]any, 0, len(list))
	for _, item := range list {
		value := openapiLatestValuePayload(item["payload"])
		if agg.legacy {
			value = openapiHistoryAggregateField(value, agg.field)
		}
		values = append(values, map[string]any{
			"value":     value,
			"quality":   openapiQualityGood,
			"timeStamp": openapiAnyInt64(item["timestamp"]),
		})
	}
	return values
}

type openapiHistoryAggregation struct {
	enabled    bool
	intervalMs int64
	function   string
	field      string
	fields     []openapiHistoryAggregationField
	legacy     bool
}

type openapiHistoryAggregationField struct {
	name     string
	function string
}

func parseOpenapiHistoryAggregation(in *types.HistoryAggregation) (openapiHistoryAggregation, error) {
	if in == nil {
		return openapiHistoryAggregation{}, nil
	}
	interval := strings.TrimSpace(in.Interval)
	function := strings.ToLower(strings.TrimSpace(in.Function))
	field := strings.TrimSpace(in.Field)
	if interval == "" {
		return openapiHistoryAggregation{}, fmt.Errorf("history aggregation requires interval")
	}
	if field != "" && len(in.Fields) > 0 {
		return openapiHistoryAggregation{}, fmt.Errorf("history aggregation field and fields cannot be used together")
	}
	dur, err := parseOpenapiHistoryInterval(interval)
	if err != nil {
		return openapiHistoryAggregation{}, err
	}
	if field != "" {
		if function == "" {
			return openapiHistoryAggregation{}, fmt.Errorf("history aggregation field requires function")
		}
		if strings.Contains(field, ".") {
			return openapiHistoryAggregation{}, fmt.Errorf("history aggregation field only supports first-level value field")
		}
		if err := validateOpenapiHistoryFunction(function); err != nil {
			return openapiHistoryAggregation{}, err
		}
		return openapiHistoryAggregation{enabled: true, intervalMs: dur.Milliseconds(), function: function, field: field, legacy: true}, nil
	}
	if function != "" {
		return openapiHistoryAggregation{}, fmt.Errorf("history aggregation function requires field")
	}
	fields := make([]openapiHistoryAggregationField, 0, len(in.Fields))
	seen := make(map[string]struct{}, len(in.Fields))
	for _, item := range in.Fields {
		name := strings.TrimSpace(item.Name)
		if name == "" || strings.Contains(name, ".") {
			return openapiHistoryAggregation{}, fmt.Errorf("history aggregation field invalid")
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			return openapiHistoryAggregation{}, fmt.Errorf("history aggregation field duplicated: %s", name)
		}
		seen[key] = struct{}{}
		itemFunction := strings.ToLower(strings.TrimSpace(item.Function))
		if itemFunction != "" {
			if err := validateOpenapiHistoryFunction(itemFunction); err != nil {
				return openapiHistoryAggregation{}, err
			}
		}
		fields = append(fields, openapiHistoryAggregationField{name: name, function: itemFunction})
	}
	return openapiHistoryAggregation{enabled: true, intervalMs: dur.Milliseconds(), fields: fields}, nil
}

func validateOpenapiHistoryFunction(function string) error {
	switch function {
	case "avg", "min", "max", "sum", "count", "first", "last":
		return nil
	default:
		return fmt.Errorf("unsupported history aggregation function: %s", function)
	}
}

func aggregateOpenapiHistory(values []map[string]any, agg openapiHistoryAggregation) ([]map[string]any, error) {
	buckets := map[int64][]map[string]any{}
	keys := make([]int64, 0)
	for _, value := range values {
		ts := openapiAnyInt64(value["timeStamp"])
		bucket := (ts / agg.intervalMs) * agg.intervalMs
		if _, ok := buckets[bucket]; !ok {
			keys = append(keys, bucket)
		}
		buckets[bucket] = append(buckets[bucket], value)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	out := make([]map[string]any, 0, len(keys))
	for _, key := range keys {
		value, err := aggregateOpenapiHistoryBucket(buckets[key], agg)
		if err != nil {
			return nil, err
		}
		out = append(out, map[string]any{
			"value":     value,
			"quality":   aggregateOpenapiQuality(buckets[key]),
			"timeStamp": key,
		})
	}
	return out, nil
}

func aggregateOpenapiHistoryBucket(values []map[string]any, agg openapiHistoryAggregation) (any, error) {
	if len(values) == 0 {
		return nil, nil
	}
	switch agg.function {
	case "first":
		return openapiHistoryAggregateField(values[0]["value"], agg.field), nil
	case "last":
		return openapiHistoryAggregateField(values[len(values)-1]["value"], agg.field), nil
	}
	nums := make([]float64, 0, len(values))
	for _, value := range values {
		raw := openapiHistoryAggregateField(value["value"], agg.field)
		num, ok := openapiFloat(raw)
		if !ok {
			if agg.function == "count" {
				continue
			}
			return nil, fmt.Errorf("history aggregation field %s requires numeric value", agg.field)
		}
		nums = append(nums, num)
	}
	if agg.function == "count" {
		return len(nums), nil
	}
	if len(nums) == 0 {
		return nil, nil
	}
	switch agg.function {
	case "sum":
		sum := 0.0
		for _, num := range nums {
			sum += num
		}
		return sum, nil
	case "avg":
		sum := 0.0
		for _, num := range nums {
			sum += num
		}
		return sum / float64(len(nums)), nil
	case "min":
		ret := math.Inf(1)
		for _, num := range nums {
			if num < ret {
				ret = num
			}
		}
		return ret, nil
	case "max":
		ret := math.Inf(-1)
		for _, num := range nums {
			if num > ret {
				ret = num
			}
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("unsupported history aggregation function: %s", agg.function)
	}
}

func openapiHistoryAggregateField(value any, field string) any {
	if body, ok := value.(map[string]any); ok {
		return body[field]
	}
	if field == "value" {
		return value
	}
	return nil
}

func aggregateOpenapiQuality(values []map[string]any) string {
	worst := openapiQualityGood
	for _, value := range values {
		switch strings.TrimSpace(fmt.Sprint(value["quality"])) {
		case "Bad":
			return "Bad"
		case "Uncertain":
			worst = "Uncertain"
		}
	}
	return worst
}

func findOpenapiNode(ctx context.Context, topic string, includeDeleted bool) (repo.UnsNode, bool, error) {
	topic = normalizeOpenapiTopic(topic)
	if topic == "" {
		return repo.UnsNode{}, false, nil
	}
	unsRepo := repo.NewUnsRepo(ctx)
	if id, err := strconv.ParseInt(topic, 10, 64); err == nil && id > 0 {
		node, err := unsRepo.GetUnsNode(ctx, id, includeDeleted)
		if errors.Is(err, repo.ErrNotFound) {
			return repo.UnsNode{}, false, nil
		}
		return node, err == nil, err
	}
	var node repo.UnsNode
	var err error
	if includeDeleted {
		node, err = unsRepo.GetUnsNodeByNamespaceIncludeDeleted(ctx, topic)
	} else {
		node, err = unsRepo.GetUnsNodeByNamespace(ctx, topic)
	}
	if errors.Is(err, repo.ErrNotFound) {
		if includeDeleted {
			node, err = unsRepo.GetUnsNodeByAliasIncludeDeleted(ctx, topic)
		} else {
			node, err = unsRepo.GetUnsNodeByAlias(ctx, topic)
		}
	}
	if errors.Is(err, repo.ErrNotFound) {
		return repo.UnsNode{}, false, nil
	}
	return node, err == nil, err
}

func expandOpenapiTopics(ctx context.Context, topics []string) ([]string, error) {
	out := make([]string, 0, len(topics))
	seen := map[string]struct{}{}
	var allNodes []repo.UnsNode
	for _, raw := range topics {
		topic := normalizeOpenapiTopic(raw)
		if topic == "" {
			continue
		}
		if !hasOpenapiWildcard(topic) {
			if _, ok := seen[topic]; !ok {
				seen[topic] = struct{}{}
				out = append(out, topic)
			}
			continue
		}
		if allNodes == nil {
			nodes, err := repo.NewUnsRepo(ctx).ListUnsNodes(ctx, repo.UnsNodeFilter{})
			if err != nil {
				return nil, err
			}
			allNodes = nodes
		}
		for _, node := range allNodes {
			if node.Type != 2 || !mqttTopicMatches(topic, node.Namespace) {
				continue
			}
			if _, ok := seen[node.Namespace]; ok {
				continue
			}
			seen[node.Namespace] = struct{}{}
			out = append(out, node.Namespace)
		}
	}
	return out, nil
}

func mqttTopicMatches(pattern, topic string) bool {
	patternParts := strings.Split(pattern, "/")
	topicParts := strings.Split(topic, "/")
	for i, part := range patternParts {
		if part == "#" {
			return true
		}
		if i >= len(topicParts) {
			return false
		}
		if part != "+" && part != topicParts[i] {
			return false
		}
	}
	return len(patternParts) == len(topicParts)
}

func normalizeOpenapiTopic(topic string) string {
	return strings.Trim(strings.TrimSpace(topic), "/")
}

func hasOpenapiWildcard(topic string) bool {
	return strings.ContainsAny(topic, "+#")
}

func openapiItemError(code int64, message string) map[string]any {
	return map[string]any{"code": code, "message": message}
}

func openapiExternalNodeType(value int16) string {
	if value == 2 {
		return "TOPIC"
	}
	return "PATH"
}

func openapiInternalNodeType(value int16) string {
	if value == 2 {
		return "file"
	}
	return "folder"
}

func openapiExternalTopicType(value int16) string {
	switch value {
	case 1:
		return "STATE"
	case 2:
		return "ACTION"
	case 3:
		return "METRIC"
	default:
		return ""
	}
}

func openapiInternalTopicType(value int16) string {
	switch value {
	case 1:
		return "State"
	case 2:
		return "Action"
	case 3:
		return "Metric"
	default:
		return ""
	}
}

func openapiTopicTypeValue(value string) (int16, bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return 0, false, nil
	case "state", "1":
		return 1, true, nil
	case "action", "2":
		return 2, true, nil
	case "metric", "3":
		return 3, true, nil
	default:
		return 0, false, fmt.Errorf("invalid topicType")
	}
}

func openapiTopicCategoryByName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "state":
		return "State"
	case "action":
		return "Action"
	case "metric":
		return "Metric"
	default:
		return ""
	}
}

func normalizeNamespaceNodeType(value string, hasChildren bool) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "path", "folder", "directory", "dir", "1":
		return "folder"
	case "topic", "file", "object", "2":
		return "file"
	case "":
		if hasChildren {
			return "folder"
		}
		return "file"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func openapiSchemaFields(raw json.RawMessage) []map[string]any {
	var body any
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &body)
	}
	var list []any
	switch typed := body.(type) {
	case []any:
		list = typed
	case map[string]any:
		if fields, ok := typed["fields"].([]any); ok {
			list = fields
		}
	}
	out := make([]map[string]any, 0, len(list))
	for _, item := range list {
		body, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name := strings.TrimSpace(fmt.Sprint(body["name"]))
		if name == "" {
			continue
		}
		field := map[string]any{
			"name": name,
			"type": strings.TrimSpace(fmt.Sprint(body["type"])),
		}
		if unit := strings.TrimSpace(fmt.Sprint(body["unit"])); unit != "" {
			field["unit"] = unit
		}
		out = append(out, field)
	}
	return out
}

func openapiExtendProperties(raw json.RawMessage) map[string]string {
	var body map[string]any
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &body)
	}
	out := make(map[string]string, len(body))
	for key, value := range body {
		out[key] = fmt.Sprint(value)
	}
	return out
}

func openapiSchemaFieldsJSON(fields []types.SchemaField) string {
	if len(fields) == 0 {
		return "[]"
	}
	body, err := json.Marshal(fields)
	if err != nil {
		return "[]"
	}
	return string(body)
}

func openapiStringMapJSON(value map[string]string) string {
	if len(value) == 0 {
		return "{}"
	}
	body, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(body)
}

func openapiLatestValuePayload(value any) any {
	body, ok := value.(map[string]any)
	if !ok {
		return value
	}
	if len(body) == 1 {
		if v, ok := body["value"]; ok {
			return v
		}
	}
	out := make(map[string]any, len(body))
	for key, item := range body {
		lower := strings.ToLower(strings.TrimSpace(key))
		if lower == "timestamp" || lower == "quality" || lower == "_timestamp" || lower == "_quality" || lower == "_ct" || lower == "_qos" || lower == "created_time" {
			continue
		}
		out[key] = item
	}
	return out
}

func openapiPayloadQuality(value any) string {
	switch strings.ToLower(strings.TrimSpace(fmt.Sprint(value))) {
	case strings.ToLower(openapiQualityBad):
		return openapiQualityBad
	case strings.ToLower(openapiQualityUncertain):
		return openapiQualityUncertain
	case "", "<nil>", strings.ToLower(openapiQualityGood):
		return openapiQualityGood
	default:
		return openapiQualityGood
	}
}

func openapiLatestTimestamp(payload map[string]any) int64 {
	for _, key := range []string{"timeStamp", "updateTime", "timestamp"} {
		if ts := openapiAnyInt64(payload[key]); ts > 0 {
			return ts
		}
	}
	if dt, ok := payload["dt"].(map[string]int64); ok {
		var latest int64
		for _, ts := range dt {
			if ts > latest {
				latest = ts
			}
		}
		return latest
	}
	if dt, ok := payload["dt"].(map[string]any); ok {
		var latest int64
		for _, value := range dt {
			if ts := openapiAnyInt64(value); ts > latest {
				latest = ts
			}
		}
		return latest
	}
	return 0
}

func parseOpenapiHistoryTime(name, value string) (int64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, fmt.Errorf("%s is required", name)
	}
	tm, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return 0, fmt.Errorf("%s must be RFC3339 UTC time", name)
	}
	return tm.UTC().UnixMilli(), nil
}

func parseOpenapiHistoryInterval(interval string) (time.Duration, error) {
	interval = strings.TrimSpace(interval)
	if interval == "" {
		return 0, fmt.Errorf("history interval empty")
	}
	if strings.HasSuffix(interval, "d") {
		days, err := strconv.ParseInt(strings.TrimSuffix(interval, "d"), 10, 64)
		if err != nil || days <= 0 {
			return 0, fmt.Errorf("invalid history interval: %s", interval)
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	dur, err := time.ParseDuration(interval)
	if err != nil {
		return 0, fmt.Errorf("invalid history interval: %s", interval)
	}
	if dur <= 0 {
		return 0, fmt.Errorf("history interval must be positive")
	}
	if dur%time.Millisecond != 0 {
		return 0, fmt.Errorf("history interval precision must be milliseconds or coarser")
	}
	return dur, nil
}

func normalizeOpenapiPage(page, size, defaultSize int64) (int64, int64) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = defaultSize
	}
	if size > openapiHistoryMaxSize {
		size = openapiHistoryMaxSize
	}
	return page, size
}

func pageOpenapiNodes(nodes []repo.UnsNode, page, size int64) []repo.UnsNode {
	start := (page - 1) * size
	if start >= int64(len(nodes)) {
		return []repo.UnsNode{}
	}
	end := start + size
	if end > int64(len(nodes)) {
		end = int64(len(nodes))
	}
	return nodes[start:end]
}

func pageOpenapiHistoryValues(values []map[string]any, page, size int64) []map[string]any {
	start := (page - 1) * size
	if start >= int64(len(values)) {
		return []map[string]any{}
	}
	end := start + size
	if end > int64(len(values)) {
		end = int64(len(values))
	}
	return values[start:end]
}

func openapiAnyInt64(value any) int64 {
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
		n, _ := v.Int64()
		return n
	case string:
		n, _ := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		return n
	default:
		return 0
	}
}

func openapiFloat(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case json.Number:
		n, err := v.Float64()
		return n, err == nil
	case string:
		n, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		return n, err == nil
	default:
		return 0, false
	}
}

func firstNonBlank(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}
