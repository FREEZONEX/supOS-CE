// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateModelsForNodeRedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量创建文件夹和文件(node-red导入专用)
func NewCreateModelsForNodeRedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateModelsForNodeRedLogic {
	return &CreateModelsForNodeRedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateModelsForNodeRedLogic) CreateModelsForNodeRed(req *types.CreateUnsNodeRedDto) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
