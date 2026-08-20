// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiFlowLegacyNodesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiFlowLegacyNodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiFlowLegacyNodesLogic {
	return &OpenapiFlowLegacyNodesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiFlowLegacyNodesLogic) OpenapiFlowLegacyNodes(req *types.OpenapiFlowLegacyNodesReq) (resp *types.Envelope, err error) {
	return NewOpenapiFlowLegacyLogic(l.ctx, l.svcCtx).Nodes(req)
}
