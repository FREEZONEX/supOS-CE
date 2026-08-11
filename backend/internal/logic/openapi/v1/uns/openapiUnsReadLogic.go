// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsReadLogic {
	return &OpenapiUnsReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsReadLogic) OpenapiUnsRead(req *types.ReadReq) (resp *types.Envelope, err error) {
	return openapiUnsRead(l.ctx, l.svcCtx, req)
}
