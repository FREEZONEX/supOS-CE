package route

import (
	"backend/internal/adapters/kong/handler"
	"backend/internal/adapters/kong/logic"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, kongLogic *logic.KongLogic, menuLogic *logic.MenuLogic, nodeRedHost, nodeRedPort string) {
	// Kong 路由管理接口
	kongHandler := handler.NewKongHandler(kongLogic)
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/inter-api/supos/kong/routeList",
				Handler: kongHandler.RouteListHandler(),
			},
		},
	)

	// 菜单管理接口
	menuHandler := handler.NewMenuHandler(menuLogic)
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/open-api/menu",
				Handler: menuHandler.SaveMenuHandler(),
			},
		},
	)

	// NodeRed 代理接口
	noderedHandler := handler.NewNodeRedHandler(nodeRedHost, nodeRedPort)
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/test/nodered",
				Handler: noderedHandler.ProxyNodeRedFlowsHandler(),
			},
			{
				Method:  http.MethodGet,
				Path:    "/flows/test/nodered",
				Handler: noderedHandler.ProxyNodeRedFlowsHandler(),
			},
		},
	)
}
