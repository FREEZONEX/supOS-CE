package kong

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RouteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取路由
func NewRouteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RouteListLogic {
	return &RouteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RouteListLogic) RouteList() error {
	// todo: add your logic here and delete this line

	return nil
}
