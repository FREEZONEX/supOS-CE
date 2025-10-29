package file

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量修改文件值
func NewBatchUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateLogic {
	return &BatchUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateLogic) BatchUpdate() error {
	// todo: add your logic here and delete this line

	return nil
}
