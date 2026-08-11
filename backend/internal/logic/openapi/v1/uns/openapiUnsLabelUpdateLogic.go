// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	coreuns "backend/internal/logic/core/uns"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsLabelUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsLabelUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsLabelUpdateLogic {
	return &OpenapiUnsLabelUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsLabelUpdateLogic) OpenapiUnsLabelUpdate(req *types.UnsLabelReq) (resp *types.Envelope, err error) {
	return coreuns.NewUnsLabelCrudLogic(l.ctx, l.svcCtx).Update(req)
}
