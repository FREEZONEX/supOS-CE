// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlowNodeTypesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowNodeTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowNodeTypesLogic {
	return &FlowNodeTypesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowNodeTypesLogic) FlowNodeTypes(req *types.NodeRedNodeTypeListReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.NodeTypes(l.ctx, req.FlowType)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
