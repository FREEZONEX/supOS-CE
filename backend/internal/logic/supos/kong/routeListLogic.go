// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package kong

import (
	"context"

	"backend/internal/adapters/kong/logic"
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
	kongLogic := logic.GetKongLogic(l.svcCtx.Config.Kong.Host, l.svcCtx.Config.Kong.Port)
	routes, err := kongLogic.RouteList(l.ctx)
	if err != nil {
		return nil, err
	}
	var data []types.SimpleRouteVO
	for _, route := range routes {
		data = append(data, types.SimpleRouteVO{
			Id:   route.ID,
			Name: route.Name,
			Url:  route.URL,
		})
	}
	resp = &types.RouteListResp{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	return
}
