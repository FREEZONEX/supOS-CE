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

type DeleteMissingNodeRedNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete one missing Node-RED node by id and location
func NewDeleteMissingNodeRedNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMissingNodeRedNodeLogic {
	return &DeleteMissingNodeRedNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMissingNodeRedNodeLogic) DeleteMissingNodeRedNode(req *types.NodeRedMissingNodeDeleteReq) (*types.NodeRedMissingNodeDeleteResult, error) {
	deleted, err := flowcommon.DeleteMissingRuntimeNode(l.ctx, l.svcCtx.SourceNodeRed, flowcommon.MissingNodeDeleteTarget{
		ID:     req.ID,
		FlowID: req.FlowID,
		Scope:  req.Scope,
	})
	if err != nil {
		return nil, err
	}
	return &types.NodeRedMissingNodeDeleteResult{Deleted: deleted}, nil
}
