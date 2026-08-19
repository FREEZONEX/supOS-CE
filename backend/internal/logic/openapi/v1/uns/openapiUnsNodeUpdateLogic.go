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

type OpenapiUnsNodeUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsNodeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsNodeUpdateLogic {
	return &OpenapiUnsNodeUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsNodeUpdateLogic) OpenapiUnsNodeUpdate(req *types.OpenapiUnsNodeSaveReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.Update(l.ctx, buildOpenapiUnsNodeSaveCommand(l.ctx, req))
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(normalizeOpenapiUnsData(data)), nil
}
