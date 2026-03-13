// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件取消标签
func NewCancelLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLabelLogic {
	return &CancelLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelLabelLogic) CancelLabel(req *types.StringArrayRequest) (resp *types.ResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
