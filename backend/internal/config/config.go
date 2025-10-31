package config

import (
	"backend/share/clients"
	"backend/share/clients/nodered"

	"gitee.com/unitedrhino/share/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database       conf.Database
	DatabaseSchema string `json:",default=supos,env=dbSchema"` //
	//Event      conf.EventConf
	CacheRedis  cache.ClusterConf
	KeycloakDSN string `json:",optional,env=KEYCLOAK_DSN,default=postgresql://postgresql:5432/keycloak" `

	OAuthKeyCloak clients.KeycloakConfig `json:",optional" `
	NodeRed       nodered.NodeRedConfig  `json:",optional" `
}
