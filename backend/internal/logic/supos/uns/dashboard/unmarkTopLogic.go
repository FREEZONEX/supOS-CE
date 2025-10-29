// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dashboard

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnmarkTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消置顶 Dashboard
func NewUnmarkTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnmarkTopLogic {
	return &UnmarkTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnmarkTopLogic) UnmarkTop(req *types.UnmarkRequest) (resp *types.JsonResult, err error) {
	// todo: add your logic here and delete this line

	return
}
