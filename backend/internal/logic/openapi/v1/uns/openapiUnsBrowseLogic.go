// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsBrowseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsBrowseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsBrowseLogic {
	return &OpenapiUnsBrowseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsBrowseLogic) OpenapiUnsBrowse(req *types.BrowseReq) (resp *types.Envelope, err error) {
	return openapiUnsBrowse(l.ctx, l.svcCtx, req)
}
