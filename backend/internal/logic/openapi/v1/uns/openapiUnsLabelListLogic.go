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

type OpenapiUnsLabelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsLabelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsLabelListLogic {
	return &OpenapiUnsLabelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsLabelListLogic) OpenapiUnsLabelList(req *types.PageReq) (resp *types.Envelope, err error) {
	return coreuns.NewUnsLabelListLogic(l.ctx, l.svcCtx).UnsLabelList(req)
}
