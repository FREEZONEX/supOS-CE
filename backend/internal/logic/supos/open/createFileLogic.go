// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建文件
func NewCreateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFileLogic {
	return &CreateFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFileLogic) CreateFile(req *types.CreateOpenApiFileDto) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
