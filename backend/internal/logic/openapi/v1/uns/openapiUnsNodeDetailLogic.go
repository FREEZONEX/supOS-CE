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

type OpenapiUnsNodeDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsNodeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsNodeDetailLogic {
	return &OpenapiUnsNodeDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsNodeDetailLogic) OpenapiUnsNodeDetail(req *types.NodeIdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.Detail(l.ctx, req.NodeId, false)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(normalizeOpenapiUnsData(data)), nil
}
