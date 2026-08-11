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

type OpenapiUnsLabelCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsLabelCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsLabelCreateLogic {
	return &OpenapiUnsLabelCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsLabelCreateLogic) OpenapiUnsLabelCreate(req *types.UnsLabelReq) (resp *types.Envelope, err error) {
	return coreuns.NewUnsLabelCrudLogic(l.ctx, l.svcCtx).Create(req)
}
