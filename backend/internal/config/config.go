package config

import (
	"backend/share/clients"

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
	KeycloakDSN string `json:",optional,env=KEYCLOAK_DSN" mapstructure:"KeycloakDSN"`

	OAuthKeyCloak clients.KeycloakConfig `json:",optional" mapstructure:"OAuthKeyCloak"`
	//DevLink    conf.DevLinkConf //和设备交互的设置
	//OssConf    conf.OssConf     `json:",optional"`
}
