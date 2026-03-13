// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改文件夹
func NewUpdateFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFolderLogic {
	return &UpdateFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFolderLogic) UpdateFolder(req *types.UpdateOpenApiFolderDto) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
