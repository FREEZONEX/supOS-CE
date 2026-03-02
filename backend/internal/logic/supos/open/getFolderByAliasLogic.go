// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFolderByAliasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 别名查询文件夹详情
func NewGetFolderByAliasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFolderByAliasLogic {
	return &GetFolderByAliasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFolderByAliasLogic) GetFolderByAlias() (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
