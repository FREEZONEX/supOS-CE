// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsDeleteLogic {
	return &OpenapiUnsDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsDeleteLogic) OpenapiUnsDelete(req *types.NodeDeleteReq) (resp *types.Envelope, err error) {
	return openapiUnsDelete(l.ctx, l.svcCtx, req)
}
