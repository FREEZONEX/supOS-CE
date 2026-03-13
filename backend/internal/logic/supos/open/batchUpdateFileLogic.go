// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量写文件实时值
func NewBatchUpdateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateFileLogic {
	return &BatchUpdateFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateFileLogic) BatchUpdateFile(req *types.UpdateFileDTOArray) (resp *types.BatchUpdateFileResult, err error) {
	// todo: add your logic here and delete this line

	return
}
