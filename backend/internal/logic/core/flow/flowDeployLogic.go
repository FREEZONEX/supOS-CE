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

type FlowDeployLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowDeployLogic {
	return &FlowDeployLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowDeployLogic) FlowDeploy(req *types.FlowDeployReq) (resp *types.Envelope, err error) {
	userID := logicx.UserID(l.ctx)
	data, err := l.svcCtx.App.Flow.Deploy(l.ctx, req.FlowId, userID, req.FlowData)
	if err != nil {
		return nil, logicx.Error(err)
	}
	recordFlowAudit(l.ctx, l.svcCtx, data, auditdomain.BizDeploy)
	return respx.Envelope(data), nil
}
