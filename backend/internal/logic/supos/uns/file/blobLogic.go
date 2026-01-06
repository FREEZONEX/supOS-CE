package file

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文件BLOB类型的值
func NewBlobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlobLogic {
	return &BlobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlobLogic) Blob() error {
	// todo: add your logic here and delete this line

	return nil
}
