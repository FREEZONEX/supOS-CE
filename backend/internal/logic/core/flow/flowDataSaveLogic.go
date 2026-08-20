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

type FlowDataSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowDataSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowDataSaveLogic {
	return &FlowDataSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowDataSaveLogic) FlowDataSave(req *types.FlowDataReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.SaveData(l.ctx, req.FlowId, req.FlowData, logicx.UserID(l.ctx))
	if err != nil {
		return nil, logicx.Error(err)
	}
	recordFlowAudit(l.ctx, l.svcCtx, data, auditdomain.BizUpdate)
	return respx.Envelope(data), nil
}
