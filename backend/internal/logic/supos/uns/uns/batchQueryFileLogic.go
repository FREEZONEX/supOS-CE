// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchQueryFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量查询文件实时值
func NewBatchQueryFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchQueryFileLogic {
	return &BatchQueryFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchQueryFileLogic) BatchQueryFile(req *types.BatchQueryFileReq) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
