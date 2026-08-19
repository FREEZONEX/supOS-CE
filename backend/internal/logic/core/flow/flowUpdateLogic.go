// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	auditdomain "backend/internal/domain/audit"
	domainflow "backend/internal/domain/flow"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlowUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowUpdateLogic {
	return &FlowUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowUpdateLogic) FlowUpdate(req *types.FlowSaveReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.Update(l.ctx, domainflow.SaveCommand{
		ID: req.FlowId, ParentID: req.ParentId, FlowType: req.FlowType, NodeType: req.NodeType, Name: req.Name,
		Description: req.Description, Template: req.Template, UnsNodeIDs: req.UnsNodeIds, UserID: logicx.UserID(l.ctx),
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	recordFlowAudit(l.ctx, l.svcCtx, data, auditdomain.BizUpdate)
	return respx.Envelope(data), nil
}
