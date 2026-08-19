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

type FlowMissingNodeCleanupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowMissingNodeCleanupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowMissingNodeCleanupLogic {
	return &FlowMissingNodeCleanupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowMissingNodeCleanupLogic) FlowMissingNodeCleanup(req *types.NodeRedMissingNodeCleanupReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.CleanupMissingNodes(l.ctx, req.FlowType)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
