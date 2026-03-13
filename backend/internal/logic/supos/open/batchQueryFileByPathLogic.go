// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchQueryFileByPathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据文件路径批量查询文件实时值
func NewBatchQueryFileByPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchQueryFileByPathLogic {
	return &BatchQueryFileByPathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchQueryFileByPathLogic) BatchQueryFileByPath(req *types.StringArrayRequest) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
