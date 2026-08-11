// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsCreateLogic {
	return &OpenapiUnsCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsCreateLogic) OpenapiUnsCreate(req *types.NodeCreateReq) (resp *types.Envelope, err error) {
	return openapiUnsCreate(l.ctx, l.svcCtx, req)
}
