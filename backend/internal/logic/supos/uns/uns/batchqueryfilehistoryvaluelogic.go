// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchQueryFileHistoryValueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量查询文件历史值
func NewBatchQueryFileHistoryValueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchQueryFileHistoryValueLogic {
	return &BatchQueryFileHistoryValueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchQueryFileHistoryValueLogic) BatchQueryFileHistoryValue(req *types.HistoryValueRequest) (resp *types.UnsHistoryQueryResult, err error) {
	// todo: add your logic here and delete this line

	return
}
