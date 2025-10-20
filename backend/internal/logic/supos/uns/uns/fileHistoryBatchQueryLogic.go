package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileHistoryBatchQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量查询文件历史值
func NewFileHistoryBatchQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileHistoryBatchQueryLogic {
	return &FileHistoryBatchQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileHistoryBatchQueryLogic) FileHistoryBatchQuery() error {
	// todo: add your logic here and delete this line

	return nil
}
