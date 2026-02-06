package handler

import (
	"backend/internal/handler/supos/open"
	imexport "backend/internal/handler/supos/uns/importExport"
	"backend/internal/handler/supos/uns/system"
	unsHandler "backend/internal/handler/supos/uns/uns"
	"backend/internal/svc"
	"net/http"
	"net/http/pprof"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterExtHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	addProfiles(server)
	addSwaggerRoutes(server)

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
			Handler: imexport.ExportHandler,
		}, rest.Route{
			Method:  http.MethodPost,
			Path:    "/inter-api/supos/uns/importExport/export/global",
			Handler: imexport.ExportGlobalHandler,
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
func addSwaggerRoutes(server *rest.Server) {
	// 提供 swagger.yaml 文件
	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/swagger/Tier0.openapi.yaml",
			Handler: open.SwaggerYamlHandler,
		},
	}, rest.WithTimeout(0))
}

func addProfiles(server *rest.Server) {
	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/",
			Handler: pprof.Index,
		},
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/allocs",
			Handler: pprof.Index,
		},
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/heap",
			Handler: pprof.Index,
		},
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/mutex",
			Handler: pprof.Index,
		},
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/threadcreate",
			Handler: pprof.Index,
		},
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/goroutine",
			Handler: pprof.Index,
		},
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/block",
			Handler: pprof.Index,
		},
		//
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/cmdline",
			Handler: pprof.Cmdline,
		},
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/profile",
			Handler: pprof.Profile,
		},
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/symbol",
			Handler: pprof.Symbol,
		},
		{
			Method:  http.MethodGet,
			Path:    "/debug/pprof/trace",
			Handler: pprof.Trace,
		},
	}, rest.WithTimeout(0))
}
