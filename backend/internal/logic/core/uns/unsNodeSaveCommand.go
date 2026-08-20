package uns

import (
	"context"

	domainuns "backend/internal/domain/uns"
	"backend/internal/logic/logicx"
	"backend/internal/types"
)

func BuildUnsNodeSaveCommand(ctx context.Context, req *types.UnsNodeSaveReq) domainuns.SaveCommand {
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
		Schema:           req.Schema,
		ExtendProperties: req.ExtendProperties,
		LabelIDs:         req.LabelIds,
		AssetFileIDs:     req.AssetFileIds,
		EnableHistory:    req.Persistence,
		AddFlow:          req.WithFlow || req.AddFlow || req.MockData,
		MockData:         req.MockData || req.AddFlow || req.WithFlow,
		UserID:           logicx.UserID(ctx),
	}
}
