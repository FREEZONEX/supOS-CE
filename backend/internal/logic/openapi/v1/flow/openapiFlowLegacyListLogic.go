// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiFlowLegacyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiFlowLegacyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiFlowLegacyListLogic {
	return &OpenapiFlowLegacyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiFlowLegacyListLogic) OpenapiFlowLegacyList(req *types.OpenapiFlowLegacyListReq) (resp *types.Envelope, err error) {
	return NewOpenapiFlowLegacyLogic(l.ctx, l.svcCtx).List(req)
}
