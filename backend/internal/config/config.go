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
	OssConf        conf.OssConf `json:",optional"`
	Export         ExportConfig `json:"export,optional"`
	GrafanaUrl     string       `json:"grafanaUrl,optional,default=http://grafana:3000"`
	PostgresqlUrl  string       `json:"sink.postgresql,env=SINK_PG_URL,optional,default=postgres://postgres:postgres@tsdb:5432/postgres"`
	TimescaledbUrl string       `json:"sink.timescaledb,env=SINK_TSDB_URL,optional,default=postgres://postgres:postgres@tsdb:5432/postgres"`
	DevLink        conf.EventConf
	CacheRedis     cache.ClusterConf
	KeycloakDSN    string                 `json:",optional,env=KEYCLOAK_DSN,default=postgresql://postgresql:5432/keycloak" `
	OAuthKeyCloak  clients.KeycloakConfig `json:",optional" `
	NodeRed        nodered.NodeRedConfig  `json:",optional" `
	Kong           clients.KongConfig     `json:",optional" mapstructure:"Kong"`
}

// ExportConfig UNS 导入导出配置
type ExportConfig struct {
	BuffeSize int `json:"buffeSize,optional,default=4096"`
	BatchSize int `json:"batchSize,optional,default=1000"`
	// 导出认定为小文件的最大数据行数
	LimitSmallFileRows int64 `json:"limitSmallFileRows,optional,default=10000"`
}

// ElasticsearchConfig represents Elasticsearch adapter configuration
type ElasticsearchConfig struct {
	Enabled   bool     `json:"enabled,optional"`
	Addresses []string `json:"addresses,optional"`
	Username  string   `json:"username,optional"`
	Password  string   `json:"password,optional,env=ELASTICSEARCH_PASSWORD"`
}
