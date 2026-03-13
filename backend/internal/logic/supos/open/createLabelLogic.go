// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建标签
func NewCreateLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLabelLogic {
	return &CreateLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLabelLogic) CreateLabel(req *types.CreateLabelDto) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
