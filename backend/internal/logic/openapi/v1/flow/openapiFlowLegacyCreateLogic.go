// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiFlowLegacyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiFlowLegacyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiFlowLegacyCreateLogic {
	return &OpenapiFlowLegacyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiFlowLegacyCreateLogic) OpenapiFlowLegacyCreate(req *types.OpenapiFlowLegacyCreateReq) (resp *types.Envelope, err error) {
	return NewOpenapiFlowLegacyLogic(l.ctx, l.svcCtx).Create(req)
}
