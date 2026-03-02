// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileByPathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 路径查询文件详情
func NewGetFileByPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileByPathLogic {
	return &GetFileByPathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileByPathLogic) GetFileByPath(req *types.GetByPathReq) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
