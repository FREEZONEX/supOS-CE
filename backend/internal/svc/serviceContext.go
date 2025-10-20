package svc

import (
	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/repo/relationDB"

	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	InitCtxsWare   rest.Middleware
	CheckTokenWare rest.Middleware
	SnowFlake      *utils.SnowFlake
	Redis          *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	stores.InitConn(c.Database)
	caches.InitStore(c.CacheRedis)
	nodeID := utils.GetNodeID(c.CacheRedis, c.Name)
	relationDB.Migrate(c.Database, c.DatabaseSchema)
	return &ServiceContext{
		Config:         c,
		CheckTokenWare: middleware.NewCheckTokenWareMiddleware().Handle,
		InitCtxsWare:   middleware.NewInitCtxsWareMiddleware().Handle,
		Redis:          redis.MustNewRedis(c.CacheRedis[0].RedisConf),
		SnowFlake:      utils.NewSnowFlake(nodeID),
	}
}
