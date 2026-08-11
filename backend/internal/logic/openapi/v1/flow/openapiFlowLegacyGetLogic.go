// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiFlowLegacyGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiFlowLegacyGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiFlowLegacyGetLogic {
	return &OpenapiFlowLegacyGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiFlowLegacyGetLogic) OpenapiFlowLegacyGet(req *types.OpenapiFlowLegacyGetReq) (resp *types.Envelope, err error) {
	return NewOpenapiFlowLegacyLogic(l.ctx, l.svcCtx).Get(req)
}
