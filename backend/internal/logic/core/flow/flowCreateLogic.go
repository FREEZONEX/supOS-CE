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

type FlowCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowCreateLogic {
	return &FlowCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowCreateLogic) FlowCreate(req *types.FlowSaveReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.Create(l.ctx, domainflow.SaveCommand{
		ParentID: req.ParentId, FlowType: req.FlowType, NodeType: req.NodeType, Name: req.Name,
		Description: req.Description, Template: req.Template, UnsNodeIDs: req.UnsNodeIds, UserID: logicx.UserID(l.ctx),
		MockData: req.MockData, MockTopic: req.MockTopic, MockFields: flowMockFieldsFromReq(req.MockFields),
		MockTriggerMode: req.MockTriggerMode,
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	recordFlowAudit(l.ctx, l.svcCtx, data, auditdomain.BizCreate)
	return respx.Envelope(data), nil
}

func flowMockFieldsFromReq(items []types.FlowMockField) []domainflow.MockField {
	fields := make([]domainflow.MockField, 0, len(items))
	for _, item := range items {
		fields = append(fields, domainflow.MockField{
			Name: item.Name,
			Type: item.FieldType,
		})
	}
	return fields
}
