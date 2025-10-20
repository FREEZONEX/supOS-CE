package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除指定路径下的所有文件夹和文件
func NewDeletePathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePathLogic {
	return &DeletePathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePathLogic) DeletePath() error {
	// todo: add your logic here and delete this line

	return nil
}
