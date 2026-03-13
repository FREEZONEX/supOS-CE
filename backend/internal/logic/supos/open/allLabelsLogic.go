// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllLabelsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标签列表
func NewAllLabelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLabelsLogic {
	return &AllLabelsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLabelsLogic) AllLabels(req *types.LabelQueryDto) (resp *types.LabelPageResult, err error) {
	// todo: add your logic here and delete this line

	return
}
