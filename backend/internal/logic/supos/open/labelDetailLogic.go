// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LabelDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标签详情
func NewLabelDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LabelDetailLogic {
	return &LabelDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LabelDetailLogic) LabelDetail() (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
