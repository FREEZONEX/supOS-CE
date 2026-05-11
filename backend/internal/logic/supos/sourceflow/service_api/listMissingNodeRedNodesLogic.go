// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package service_api

import (
	"context"

	"backend/internal/logic/supos/flowcommon"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMissingNodeRedNodesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List missing Node-RED nodes across all source flow tabs
func NewListMissingNodeRedNodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMissingNodeRedNodesLogic {
	return &ListMissingNodeRedNodesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMissingNodeRedNodesLogic) ListMissingNodeRedNodes() (*types.NodeRedMissingNodeListResult, error) {
	nodes, err := flowcommon.ListMissingRuntimeNodes(l.ctx, l.svcCtx.SourceNodeRed)
	if err != nil {
		return nil, err
	}
	resp := &types.NodeRedMissingNodeListResult{
		Nodes: make([]types.NodeRedMissingNode, 0, len(nodes)),
	}
	for _, node := range nodes {
		resp.Nodes = append(resp.Nodes, types.NodeRedMissingNode{
			ID:        node.ID,
			Type:      node.Type,
			Name:      node.Name,
			Scope:     node.Scope,
			FlowID:    node.FlowID,
			FlowLabel: node.FlowLabel,
			Users:     node.Users,
		})
	}
	return resp, nil
}
