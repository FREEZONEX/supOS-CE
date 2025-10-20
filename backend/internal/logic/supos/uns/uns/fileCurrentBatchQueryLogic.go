package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileCurrentBatchQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量查询文件实时值
func NewFileCurrentBatchQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileCurrentBatchQueryLogic {
	return &FileCurrentBatchQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileCurrentBatchQueryLogic) FileCurrentBatchQuery() error {
	// todo: add your logic here and delete this line

	return nil
}
