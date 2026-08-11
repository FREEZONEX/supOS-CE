// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsSearchLogic {
	return &OpenapiUnsSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsSearchLogic) OpenapiUnsSearch(req *types.SearchReq) (resp *types.Envelope, err error) {
	return openapiUnsSearch(l.ctx, l.svcCtx, req)
}
