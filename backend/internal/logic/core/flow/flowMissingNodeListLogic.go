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

type FlowMissingNodeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowMissingNodeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowMissingNodeListLogic {
	return &FlowMissingNodeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowMissingNodeListLogic) FlowMissingNodeList(req *types.NodeRedMissingNodeListReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.ListMissingNodes(l.ctx, req.FlowType)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
