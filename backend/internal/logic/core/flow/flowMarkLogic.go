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

type FlowMarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowMarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowMarkLogic {
	return &FlowMarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowMarkLogic) FlowMark(req *types.FlowMarkReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.Mark(l.ctx, req.FlowId, req.Pinned, logicx.UserID(l.ctx))
	if err != nil {
		return nil, logicx.Error(err)
	}
	businessType := auditdomain.BizUnmark
	if req.Pinned {
		businessType = auditdomain.BizMark
	}
	recordFlowAudit(l.ctx, l.svcCtx, data, businessType)
	return respx.Envelope(data), nil
}
