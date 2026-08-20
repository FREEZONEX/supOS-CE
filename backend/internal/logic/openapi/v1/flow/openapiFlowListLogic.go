// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"context"

	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiFlowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiFlowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiFlowListLogic {
	return &OpenapiFlowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiFlowListLogic) OpenapiFlowList(req *types.FlowListReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.Flow.List(l.ctx, req.FlowType, req.ParentId, req.Keyword)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
