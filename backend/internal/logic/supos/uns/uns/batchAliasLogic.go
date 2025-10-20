package uns

import (
	"context"

	"backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type BatchAliasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据别名集合批量删除文件夹和文件
func NewBatchAliasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchAliasLogic {
	return &BatchAliasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchAliasLogic) BatchAlias() error {
	// todo: add your logic here and delete this line

	return nil
}
