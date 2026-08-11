// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package svc

import (
	"backend/internal/bootstrap"
	"backend/internal/config"
	"backend/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config      config.Config
	App         *bootstrap.App
	SessionAuth rest.Middleware
	Permission  rest.Middleware
	ApiKeyAuth  rest.Middleware
}

func NewServiceContext(c config.Config, app *bootstrap.App) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		App:         app,
		SessionAuth: middleware.NewSessionAuthMiddleware(app.Auth, app.Audit).Handle,
		Permission:  middleware.NewPermissionMiddleware(app.Permission).Handle,
		ApiKeyAuth:  middleware.NewApiKeyAuthMiddleware(app.Auth, app.Permission, app.Audit).Handle,
	}
}
