// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改标签
func NewUpdateLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLabelLogic {
	return &UpdateLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLabelLogic) UpdateLabel(req *types.UpdateLabelDto) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
