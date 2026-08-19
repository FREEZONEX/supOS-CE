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

type OpenapiUnsLabelDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsLabelDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsLabelDetailLogic {
	return &OpenapiUnsLabelDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsLabelDetailLogic) OpenapiUnsLabelDetail(req *types.IdReq) (resp *types.Envelope, err error) {
	return coreuns.NewUnsLabelCrudLogic(l.ctx, l.svcCtx).Detail(req)
}
