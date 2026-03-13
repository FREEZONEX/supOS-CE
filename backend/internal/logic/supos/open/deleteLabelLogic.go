// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除标签
func NewDeleteLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLabelLogic {
	return &DeleteLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLabelLogic) DeleteLabel() (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
