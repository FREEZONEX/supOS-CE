// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package kong

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RouteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取简化的路由列表
func NewRouteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RouteListLogic {
	return &RouteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RouteListLogic) RouteList() (resp *types.RouteListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
