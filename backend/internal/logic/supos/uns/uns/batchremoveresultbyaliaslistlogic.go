package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchRemoveResultByAliasListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据别名集合批量删除文件夹和文件
func NewBatchRemoveResultByAliasListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchRemoveResultByAliasListLogic {
	return &BatchRemoveResultByAliasListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchRemoveResultByAliasListLogic) BatchRemoveResultByAliasList(req *types.BatchRemoveUnsDto) (resp *types.RemoveResult, err error) {
	// todo: add your logic here and delete this line

	return
}
