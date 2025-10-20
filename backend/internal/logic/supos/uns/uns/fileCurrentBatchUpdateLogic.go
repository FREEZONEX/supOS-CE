package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileCurrentBatchUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量写文件实时值
func NewFileCurrentBatchUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileCurrentBatchUpdateLogic {
	return &FileCurrentBatchUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileCurrentBatchUpdateLogic) FileCurrentBatchUpdate() error {
	// todo: add your logic here and delete this line

	return nil
}
