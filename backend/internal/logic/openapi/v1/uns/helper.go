package uns

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domainuns "backend/internal/domain/uns"
	"backend/internal/logic/logicx"
	"backend/internal/repo"
	"backend/internal/types"
)

func buildOpenapiUnsNodeSaveCommand(ctx context.Context, req *types.OpenapiUnsNodeSaveReq) domainuns.SaveCommand {
	return domainuns.SaveCommand{
		ID:               req.NodeId,
		ParentID:         req.ParentId,
		Name:             req.Name,
		DisplayName:      req.DisplayName,
		Description:      req.Description,
		Alias:            req.Alias,
		Namespace:        req.Namespace,
		NodeType:         req.NodeType,
		TopicType:        req.TopicType,
		Schema:           openapiSchemaFieldsJSON(req.Fields),
		ExtendProperties: req.ExtendProperties,
		LabelIDs:         req.LabelIds,
		AssetFileIDs:     req.AssetFileIds,
		EnableHistory:    req.EnableHistory,
		AddFlow:          req.WithFlow || req.AddFlow || req.MockData,
		MockData:         req.MockData || req.AddFlow || req.WithFlow,
		UserID:           logicx.UserID(ctx),
	}
}

func resolveOpenapiAttachmentNode(ctx context.Context, unsID int64, topic string) (repo.UnsNode, error) {
	unsRepo := repo.NewUnsRepo(ctx)
	topic = normalizeOpenapiTopic(topic)
	if topic != "" {
		node, err := unsRepo.GetUnsNodeByNamespace(ctx, topic)
		if errors.Is(err, repo.ErrNotFound) {
			node, err = unsRepo.GetUnsNodeByAlias(ctx, topic)
		}
		return node, err
	}
	if unsID <= 0 {
		return repo.UnsNode{}, fmt.Errorf("topic or unsId is required")
	}
	return unsRepo.GetUnsNode(ctx, unsID, false)
}

func normalizeOpenapiUnsData(value any) any {
	switch v := value.(type) {
	case []map[string]any:
		for i := range v {
			v[i] = normalizeOpenapiUnsData(v[i]).(map[string]any)
		}
		return v
	case []any:
		for i := range v {
			v[i] = normalizeOpenapiUnsData(v[i])
		}
		return v
	case map[string]any:
		for key, item := range v {
			v[key] = normalizeOpenapiUnsData(item)
		}
		normalizeOpenapiUnsNode(v)
		return v
	default:
		return value
	}
}

func normalizeOpenapiUnsNode(item map[string]any) {
	if !isOpenapiUnsNode(item) {
		return
	}
	if nodeType, ok := item["type"].(string); ok {
		switch strings.ToLower(strings.TrimSpace(nodeType)) {
		case "folder", "path", "1":
			item["type"] = "PATH"
		case "file", "topic", "2":
			item["type"] = "TOPIC"
		}
	}
	if topicType, ok := item["topicType"].(string); ok {
		item["topicType"] = strings.ToUpper(strings.TrimSpace(topicType))
	}
	if enableHistory, ok := boolFromOpenapiValue(item["enableHistory"]); ok {
		item["enableHistory"] = enableHistory
	} else if persistence, ok := boolFromOpenapiValue(item["persistence"]); ok {
		item["enableHistory"] = persistence
	}
	delete(item, "persistence")
}

func isOpenapiUnsNode(item map[string]any) bool {
	if _, ok := item["id"]; !ok {
		return false
	}
	if _, ok := item["type"]; !ok {
		return false
	}
	if _, ok := item["name"]; !ok {
		return false
	}
	return true
}

func boolFromOpenapiValue(value any) (bool, bool) {
	switch v := value.(type) {
	case bool:
		return v, true
	case int:
		return v == 1, true
	case int8:
		return v == 1, true
	case int16:
		return v == 1, true
	case int32:
		return v == 1, true
	case int64:
		return v == 1, true
	case uint:
		return v == 1, true
	case uint8:
		return v == 1, true
	case uint16:
		return v == 1, true
	case uint32:
		return v == 1, true
	case uint64:
		return v == 1, true
	case float32:
		return v == 1, true
	case float64:
		return v == 1, true
	case string:
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "1", "true", "yes", "y", "enable", "enabled":
			return true, true
		case "0", "2", "false", "no", "n", "disable", "disabled":
			return false, true
		}
		return false, false
	default:
		if value == nil {
			return false, false
		}
		return boolFromOpenapiValue(fmt.Sprint(value))
	}
}
