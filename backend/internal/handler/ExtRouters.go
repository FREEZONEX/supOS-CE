package handler

import (
	imexport "backend/internal/handler/supos/uns/importExport"
	"backend/internal/handler/supos/uns/system"
	unsHandler "backend/internal/handler/supos/uns/uns"
	"backend/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterExtHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {

	server.AddRoutes(rest.WithMiddlewares(
		[]rest.Middleware{serverCtx.CheckTokenWare, serverCtx.InitCtxsWare},
		rest.Route{
			Method:  http.MethodPost,
			Path:    "/inter-api/supos/uns/importExport/import",
			Handler: imexport.ImportHandler,
		},
	), rest.WithTimeout(0), rest.WithMaxBytes(1<<30))

	server.AddRoutes(rest.WithMiddlewares(
		[]rest.Middleware{serverCtx.CheckTokenWare, serverCtx.InitCtxsWare},
		rest.Route{
			Method:  http.MethodPost,
			Path:    "/inter-api/supos/uns/importExport/export",
			Handler: imexport.ExportHandler(serverCtx),
		}, rest.Route{
			Method:  http.MethodDelete, // 删除指定路径下的所有文件夹和文件，不要带超时时间
			Path:    "/inter-api/supos/uns",
			Handler: unsHandler.RemoveModelOrInstanceHandler(serverCtx),
		},
	), rest.WithTimeout(0))

	server.AddRoutes(rest.WithMiddlewares(
		[]rest.Middleware{serverCtx.CheckTokenWare, serverCtx.InitCtxsWare},
		rest.Route{
			Method:  http.MethodGet,
			Path:    "/inter-api/supos/uns/newMsg",
			Handler: unsHandler.PushNewMsgHandler,
		},
	), rest.WithTimeout(0), rest.WithSSE())

	server.AddRoutes(rest.WithMiddlewares(
		[]rest.Middleware{serverCtx.CheckTokenWare, serverCtx.InitCtxsWare},
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/inter-api/supos/uns/dev",
				Handler: system.DevtestHandler,
			}, {
				Method:  http.MethodPost,
				Path:    "/inter-api/supos/uns/dev",
				Handler: system.DevtestHandler,
			},
		}...,
	))
}
