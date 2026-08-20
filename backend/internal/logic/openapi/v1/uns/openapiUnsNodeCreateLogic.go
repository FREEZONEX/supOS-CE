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

type OpenapiUnsNodeCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsNodeCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsNodeCreateLogic {
	return &OpenapiUnsNodeCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsNodeCreateLogic) OpenapiUnsNodeCreate(req *types.OpenapiUnsNodeSaveReq) (resp *types.Envelope, err error) {
	cmd, err := l.svcCtx.App.UNS.NormalizeOpenapiCreateCommand(l.ctx, buildOpenapiUnsNodeSaveCommand(l.ctx, req))
	if err != nil {
		return nil, logicx.Error(err)
	}
	data, err := l.svcCtx.App.UNS.Create(l.ctx, cmd)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(normalizeOpenapiUnsData(data)), nil
}
