// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 置顶 Dashboard
func NewMarkTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkTopLogic {
	return &MarkTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkTopLogic) MarkTop(req *types.MarkTopRequest) (resp *types.JsonResult, err error) {
	// todo: add your logic here and delete this line

	return
}
