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
	CacheRedis    cache.ClusterConf
	KeycloakDSN   string                 `json:",optional,env=KEYCLOAK_DSN,default=postgresql://postgresql:5432/keycloak" `
	OAuthKeyCloak clients.KeycloakConfig `json:",optional" `
	NodeRed       nodered.NodeRedConfig  `json:",optional" `
	Kong          clients.KongConfig     `json:",optional" mapstructure:"Kong"`
	Elasticsearch ElasticsearchConfig    `json:",optional"`
}

// ElasticsearchConfig represents Elasticsearch adapter configuration
type ElasticsearchConfig struct {
	Enabled   bool     `json:"enabled,optional"`
	Addresses []string `json:"addresses,optional"`
	Username  string   `json:"username,optional"`
	Password  string   `json:"password,optional,env=ELASTICSEARCH_PASSWORD"`
}
