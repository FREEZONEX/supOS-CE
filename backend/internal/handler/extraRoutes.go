package handler

import (
	"net/http"

	coresystem "backend/internal/handler/core/system"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterExtraHandlers(server *rest.Server) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/swagger",
				Handler: coresystem.OpenapiV1SwaggerUIHandler(),
			},
			{
				Method:  http.MethodGet,
				Path:    "/swagger.json",
				Handler: coresystem.OpenapiV1SwaggerJSONHandler(),
			},
		},
		rest.WithPrefix("/openapi/v1"),
		rest.WithTimeout(0),
	)
}
