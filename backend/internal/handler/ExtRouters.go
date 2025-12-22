package handler

import (
	imexport "backend/internal/handler/supos/uns/importExport"
	unsHandler "backend/internal/handler/supos/uns/uns"
	"backend/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterExtHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/inter-api/supos/uns/importExport/import",
		Handler: imexport.ImportHandler(serverCtx),
	}, rest.WithTimeout(0), rest.WithMaxBytes(1<<30))

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/inter-api/supos/uns/importExport/export",
		Handler: imexport.ExportHandler(serverCtx),
	}, rest.WithTimeout(0))

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/inter-api/supos/uns/newMsg",
		Handler: unsHandler.PushNewMsgHandler,
	}, rest.WithSSE(), rest.WithTimeout(0))

}
