package svc

import (
	"backend/internal/common"
	"backend/internal/middleware"
	"backend/internal/repo/relationDB"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	cache "backend/internal/common/cache"
	"backend/internal/config"
	keycloakrepo "backend/internal/repo/keycloak"
	"backend/share/clients"
)

type ServiceContext struct {
	Config         config.Config
	InitCtxsWare   rest.Middleware
	CheckTokenWare rest.Middleware
	SnowFlake      *utils.SnowFlake
	Keycloak       *clients.KeycloakClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	relationDB.Migrate(c.Database, c.DatabaseSchema)

	if err := cache.InitCaches(); err != nil {
		logx.Errorf("failed to init cache: %v", err)
		panic(err)
	}

	dbConn, err := stores.GetConn(conf.Database{
		DBType: "pgsql",
		DSN:    c.KeycloakDSN,
	})
	if err != nil {
		logx.Errorf("failed to init keycloak database: %v", err)
	} else {
		if err := keycloakrepo.InitWithDB(dbConn); err != nil {
			logx.Errorf("failed to register keycloak database: %v", err)
		}
	}

	keycloakClient := clients.InitKeycloakClient(c.OAuthKeyCloak)

	common.InitSnowflake(1)
	return &ServiceContext{
		Config:         c,
		CheckTokenWare: middleware.NewCheckTokenWareMiddleware(keycloakClient, c.OAuthKeyCloak.SuposHome, c.OAuthKeyCloak.Realm).Handle,
		InitCtxsWare:   middleware.NewInitCtxsWareMiddleware().Handle,
		SnowFlake:      utils.NewSnowFlake(1),
		Keycloak:       keycloakClient,
	}
}
