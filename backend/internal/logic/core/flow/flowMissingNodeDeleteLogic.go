// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	flowdomain "backend/internal/domain/flow"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlowMissingNodeDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowMissingNodeDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowMissingNodeDeleteLogic {
	return &FlowMissingNodeDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowMissingNodeDeleteLogic) FlowMissingNodeDelete(req *types.NodeRedMissingNodeDeleteReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.DeleteMissingNode(l.ctx, flowdomain.MissingNodeDeleteCommand{
		ID:       req.Id,
		FlowID:   req.FlowId,
		Scope:    req.Scope,
		FlowType: req.FlowType,
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
