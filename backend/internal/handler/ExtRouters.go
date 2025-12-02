package handler

import (
	unsHandler "backend/internal/handler/supos/uns/uns"
	"backend/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterExtHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/inter-api/supos/uns/newMsg",
		Handler: unsHandler.PushNewMsgHandler(serverCtx),
	}, rest.WithSSE(), rest.WithTimeout(0))
}
