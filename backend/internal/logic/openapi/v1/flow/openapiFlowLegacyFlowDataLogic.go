// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiFlowLegacyFlowDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiFlowLegacyFlowDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiFlowLegacyFlowDataLogic {
	return &OpenapiFlowLegacyFlowDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiFlowLegacyFlowDataLogic) OpenapiFlowLegacyFlowData(req *types.OpenapiFlowLegacyGetReq) (resp *types.Envelope, err error) {
	return NewOpenapiFlowLegacyLogic(l.ctx, l.svcCtx).FlowData(req)
}
