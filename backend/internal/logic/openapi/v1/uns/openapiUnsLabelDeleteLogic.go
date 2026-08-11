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

type OpenapiUnsLabelDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsLabelDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsLabelDeleteLogic {
	return &OpenapiUnsLabelDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsLabelDeleteLogic) OpenapiUnsLabelDelete(req *types.IdReq) (resp *types.Envelope, err error) {
	return coreuns.NewUnsLabelCrudLogic(l.ctx, l.svcCtx).Delete(req)
}
