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

type FlowDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowDeleteLogic {
	return &FlowDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowDeleteLogic) FlowDelete(req *types.FlowIdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.Delete(l.ctx, req.FlowId, logicx.UserID(l.ctx))
	if err != nil {
		return nil, logicx.Error(err)
	}
	recordFlowAudit(l.ctx, l.svcCtx, data, auditdomain.BizDelete)
	return respx.Envelope(data), nil
}
