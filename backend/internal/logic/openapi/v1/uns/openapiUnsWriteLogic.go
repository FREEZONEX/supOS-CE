// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsWriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsWriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsWriteLogic {
	return &OpenapiUnsWriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsWriteLogic) OpenapiUnsWrite(req *types.WriteReq) (resp *types.Envelope, err error) {
	return openapiUnsWrite(l.ctx, l.svcCtx, req)
}
