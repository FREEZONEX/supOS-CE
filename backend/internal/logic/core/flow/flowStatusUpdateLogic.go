// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	auditdomain "backend/internal/domain/audit"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlowStatusUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowStatusUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowStatusUpdateLogic {
	return &FlowStatusUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowStatusUpdateLogic) FlowStatusUpdate(req *types.FlowStatusReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.UpdateStatus(l.ctx, req.FlowId, logicx.UserID(l.ctx), req.Status)
	if err != nil {
		return nil, logicx.Error(err)
	}
	recordFlowAudit(l.ctx, l.svcCtx, data, auditdomain.BizUpdate)
	return respx.Envelope(data), nil
}
