package route

import (
	"backend/internal/logic/supos/uns/dashboard/dao"
	"backend/internal/logic/supos/uns/dashboard/handler"
	"backend/internal/logic/supos/uns/dashboard/service"
	unsservice "backend/internal/logic/supos/uns/uns/service"
	"backend/share/spring"
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

// RegisterHandlers 注册 Dashboard 模块的路由
func RegisterHandlers(server *rest.Server, sqlConn sqlx.SqlConn) {
	ctx := context.Background()

	// 创建 DAO 层实例
	dashboardMapper := dao.NewDashboardMapper(sqlConn, ctx)
	dashboardRefMapper := dao.NewDashboardRefMapper(sqlConn, ctx)
	dashboardMarkMapper := dao.NewDashboardMarkedMapper(sqlConn, ctx)

	unsQueryService := spring.GetBean[*unsservice.UnsQueryService]()
	unsUpdateService := spring.GetBean[*unsservice.UnsUpdateService]()

	// 创建 Logic 层实例
	dashboardLogic := service.NewDashboardService(ctx, dashboardMapper, dashboardRefMapper, dashboardMarkMapper, unsQueryService, unsUpdateService)

	// 创建 Handler 层实例
	dashboardHandler := handler.NewDashboardHandler(dashboardLogic)

	// 注册路由
	server.AddRoutes(
		[]rest.Route{
			// 分页查询
			{
				Method:  "GET",
				Path:    "/inter-api/supos/uns/dashboard",
				Handler: dashboardHandler.PageListHandler(),
			},
			// 获取详情
			{
				Method:  "GET",
				Path:    "/inter-api/supos/uns/dashboard/detail",
				Handler: dashboardHandler.GetDetailHandler(),
			},
			// 根据 UID 获取
			{
				Method:  "GET",
				Path:    "/inter-api/supos/uns/dashboard/:uid",
				Handler: dashboardHandler.GetByUuidHandler(),
			},
			// 创建
			{
				Method:  "POST",
				Path:    "/inter-api/supos/uns/dashboard",
				Handler: dashboardHandler.CreateHandler(),
			},
			// 编辑
			{
				Method:  "PUT",
				Path:    "/inter-api/supos/uns/dashboard",
				Handler: dashboardHandler.EditHandler(),
			},
			// 删除
			{
				Method:  "DELETE",
				Path:    "/inter-api/supos/uns/dashboard/:uid",
				Handler: dashboardHandler.DeleteHandler(),
			},
			// 置顶
			{
				Method:  "POST",
				Path:    "/inter-api/supos/uns/dashboard/mark",
				Handler: dashboardHandler.MarkTopHandler(),
			},
			// 取消置顶
			{
				Method:  "DELETE",
				Path:    "/inter-api/supos/uns/dashboard/unmark",
				Handler: dashboardHandler.UnmarkTopHandler(),
			},
			// 绑定 UNS
			{
				Method:  "POST",
				Path:    "/inter-api/supos/uns/dashboard/bindUns",
				Handler: dashboardHandler.BindUnsHandler(),
			},
			// 根据 UNS 获取
			{
				Method:  "GET",
				Path:    "/inter-api/supos/uns/dashboard/getByUns",
				Handler: dashboardHandler.GetByUnsHandler(),
			},
			// 基于 UNS 创建 Grafana Dashboard
			{
				Method:  "POST",
				Path:    "/inter-api/supos/uns/dashboard/createGrafanaByUns/:alias",
				Handler: dashboardHandler.CreateGrafanaByUnsHandler(),
			},
			// 检查是否存在
			{
				Method:  "GET",
				Path:    "/inter-api/supos/uns/dashboard/isExist",
				Handler: dashboardHandler.IsExistHandler(),
			},
		},
	)
}
