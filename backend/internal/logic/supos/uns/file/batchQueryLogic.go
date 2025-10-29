package file

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量查询文件实时值
func NewBatchQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchQueryLogic {
	return &BatchQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchQueryLogic) BatchQuery() error {
	// todo: add your logic here and delete this line

	return nil
}
