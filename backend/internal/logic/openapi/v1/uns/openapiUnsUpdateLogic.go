// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsUpdateLogic {
	return &OpenapiUnsUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsUpdateLogic) OpenapiUnsUpdate(req *types.NodeUpdateReq) (resp *types.Envelope, err error) {
	return openapiUnsUpdate(l.ctx, l.svcCtx, req)
}
