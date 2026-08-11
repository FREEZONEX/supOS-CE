// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpenapiUnsHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpenapiUnsHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenapiUnsHistoryLogic {
	return &OpenapiUnsHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpenapiUnsHistoryLogic) OpenapiUnsHistory(req *types.HistoryReq) (resp *types.Envelope, err error) {
	return openapiUnsHistory(l.ctx, l.svcCtx, req)
}
