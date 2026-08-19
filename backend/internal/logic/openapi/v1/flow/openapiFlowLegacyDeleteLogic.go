// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiFlowLegacyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiFlowLegacyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiFlowLegacyDeleteLogic {
	return &OpenapiFlowLegacyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiFlowLegacyDeleteLogic) OpenapiFlowLegacyDelete(req *types.OpenapiFlowLegacyDeleteReq) (resp *types.Envelope, err error) {
	return NewOpenapiFlowLegacyLogic(l.ctx, l.svcCtx).Delete(req)
}
