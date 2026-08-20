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

type OpenapiUnsNodeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsNodeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsNodeListLogic {
	return &OpenapiUnsNodeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsNodeListLogic) OpenapiUnsNodeList(req *types.UnsListReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.List(l.ctx, req.ParentId, req.ParentIdSet, req.Keyword, req.IncludeRecycle)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(normalizeOpenapiUnsData(data)), nil
}
