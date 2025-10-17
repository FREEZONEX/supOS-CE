package svc

import (
	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/repo/relationDB"

	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	InitCtxsWare   rest.Middleware
	CheckTokenWare rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	relationDB.Migrate(c.Database, c.DatabaseSchema)
	return &ServiceContext{
		Config:         c,
		CheckTokenWare: middleware.NewCheckTokenWareMiddleware().Handle,
		InitCtxsWare:   middleware.NewInitCtxsWareMiddleware().Handle,
	}
}
