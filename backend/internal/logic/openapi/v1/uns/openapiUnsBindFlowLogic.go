// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsBindFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsBindFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsBindFlowLogic {
	return &OpenapiUnsBindFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsBindFlowLogic) OpenapiUnsBindFlow(req *types.UnsBindFlowReq) (resp *types.Envelope, err error) {
	if _, err := l.svcCtx.App.UNS.Detail(l.ctx, req.UnsId, false); err != nil {
		return nil, logicx.Error(err)
	}
	_, err = l.svcCtx.App.Flow.BindUns(l.ctx, req.FlowId, req.UnsId, logicx.UserID(l.ctx))
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(map[string]any{
		"success": true,
		"unsId":   req.UnsId,
		"flowId":  req.FlowId,
	}), nil
}
