package resource

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// batch
func NewBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchLogic {
	return &BatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchLogic) Batch() error {
	// todo: add your logic here and delete this line

	return nil
}
