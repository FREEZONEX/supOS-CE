// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package config

import (
	"strings"


	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	ProductVersion string `json:"productVersion,optional,env=PRODUCT_VERSION"`
	Database       DatabaseConf
	Redis          RedisConf
	Security       SecurityConf
	Gateway        GatewayConf
	Storage        StorageConf
	DataIngest     DataIngestConf
	Tier0Sdk       Tier0SdkConfig   `json:"tier0Sdk,optional"`
	Export         ExportConfig     `json:"export,optional"`
}

type DatabaseConf struct {
	UnsDbUrl                 string
	SinkDbUrl                string
	TimeseriesRetentionYears int `json:"timeseriesRetentionYears,optional"`
}

// TenantPGConf 租户 PG 隔离配置


type RedisConf struct {
	Addr     string
	Password string
	DB       int
}

type SecurityConf struct {
	JwtSecret            string
	InitialAdminPassword string
	PluginAPIKey         string
}

type GatewayConf struct {
	WebDir           string
	LocalFrontendDev bool
	FrontendDevUrl   string
	SourceFlowUrl     string
	EventFlowUrl      string
	EmqxUrl           string
}

type StorageConf struct {
	FileRoot string
}







type DataIngestConf struct {
	Enabled         bool
	MqttBrokers     string
	MqttClientID    string
	MqttUsername    string
	MqttPassword    string
	MqttTopic       string
	QueueSize       int
	BatchSize       int
	FlushIntervalMs int
}










// Tier0SdkConfig Tier0 SDK 环境变量配置
type Tier0SdkConfig struct {
	ApiHost  string `json:"apiHost,optional,env=TIER0_SDK_API_HOST"` // host:port，不带协议前缀
	MqttHost string `json:"mqttHost,optional,env=TIER0_SDK_MQTT_HOST,default=emqx"`
	MqttPort string `json:"mqttPort,optional,env=TIER0_SDK_MQTT_PORT,default=8083"`
}

// ExportConfig UNS 导入导出配置
type ExportConfig struct {
	BuffeSize          int   `json:"buffeSize,optional,default=4096"`
	BatchSize          int   `json:"batchSize,optional,default=1000"`
	LimitSmallFileRows int64 `json:"limitSmallFileRows,optional,default=10000"`
}

func (c *Config) Normalize() {
	if c.Name == "" {
		c.Name = "backend"
	}
	c.ProductVersion = strings.TrimSpace(getenv("PRODUCT_VERSION", c.ProductVersion))
	if c.ProductVersion == "" {
		c.ProductVersion = "dev"
	}
	if c.Host == "" {
		c.Host = "0.0.0.0"
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	c.Database.UnsDbUrl = expandEnv(c.Database.UnsDbUrl)
	c.Database.SinkDbUrl = expandEnv(c.Database.SinkDbUrl)
	c.Database.TimeseriesRetentionYears = NormalizeTimeseriesRetentionYears(c.Database.TimeseriesRetentionYears)
	c.Log.Mode = getenv("BACKEND_LOG_MODE", c.Log.Mode)
	c.Log.Encoding = getenv("BACKEND_LOG_ENCODING", c.Log.Encoding)
	c.Log.Path = expandEnv(getenv("BACKEND_LOG_PATH", c.Log.Path))
	c.Log.Level = getenv("BACKEND_LOG_LEVEL", c.Log.Level)
	c.Log.Rotation = getenv("BACKEND_LOG_ROTATION", c.Log.Rotation)
	c.Log.Compress = getenvBool("BACKEND_LOG_COMPRESS", c.Log.Compress)
	c.Log.KeepDays = getenvInt("BACKEND_LOG_KEEP_DAYS", c.Log.KeepDays)
	c.Redis.Addr = expandEnv(c.Redis.Addr)
	c.Redis.Password = expandEnv(c.Redis.Password)
	c.Security.JwtSecret = expandEnv(c.Security.JwtSecret)
	c.Security.InitialAdminPassword = expandEnv(c.Security.InitialAdminPassword)
	c.Security.PluginAPIKey = expandEnv(c.Security.PluginAPIKey)
	c.CertFile = expandEnv(c.CertFile)
	c.KeyFile = expandEnv(c.KeyFile)
	c.Gateway.WebDir = expandEnv(c.Gateway.WebDir)
	c.Gateway.FrontendDevUrl = expandEnv(c.Gateway.FrontendDevUrl)
	c.Gateway.SourceFlowUrl = expandEnv(c.Gateway.SourceFlowUrl)
	c.Gateway.EventFlowUrl = expandEnv(c.Gateway.EventFlowUrl)
	c.Gateway.EmqxUrl = expandEnv(c.Gateway.EmqxUrl)
	c.Storage.FileRoot = expandEnv(c.Storage.FileRoot)
	c.DataIngest.MqttBrokers = expandEnv(c.DataIngest.MqttBrokers)
	c.DataIngest.MqttClientID = expandEnv(c.DataIngest.MqttClientID)
	c.DataIngest.MqttUsername = expandEnv(c.DataIngest.MqttUsername)
	c.DataIngest.MqttPassword = expandEnv(c.DataIngest.MqttPassword)
	c.DataIngest.MqttTopic = expandEnv(c.DataIngest.MqttTopic)
	if c.Database.UnsDbUrl == "" {
		c.Database.UnsDbUrl = getenv("UNS_DB_URL", "postgres://postgres:postgres@localhost:5432/edge?sslmode=disable")
	}
	if c.Database.SinkDbUrl == "" {
		c.Database.SinkDbUrl = getenv("SINK_DB_URL", c.Database.UnsDbUrl)
	}
	if c.Redis.Addr == "" {
		c.Redis.Addr = getenv("REDIS_ADDR", "localhost:6379")
	}
	if c.Security.JwtSecret == "" {
		c.Security.JwtSecret = getenv("JWT_SECRET", "")
	}
	if c.Security.InitialAdminPassword == "" {
		c.Security.InitialAdminPassword = getenv("ADMIN_INITIAL_PASSWORD", "")
	}
	if c.Security.PluginAPIKey == "" {
		c.Security.PluginAPIKey = getenv("TIER0_API_KEY", "")
	}
	if c.Gateway.WebDir == "" {
		c.Gateway.WebDir = getenv("WEB_DIR", "./web")
	}
	c.Gateway.LocalFrontendDev = getenvBool("LOCAL_FRONTEND_DEV", c.Gateway.LocalFrontendDev)
	if c.Gateway.FrontendDevUrl == "" {
		c.Gateway.FrontendDevUrl = getenv("FRONTEND_DEV_PROXY_URL", getenv("FRONTEND_DEV_URL", ""))
	}
	if c.Gateway.SourceFlowUrl == "" {
		c.Gateway.SourceFlowUrl = getenv("SOURCEFLOW_URL", "http://sourceflow:1880")
	}
	if c.Gateway.EventFlowUrl == "" {
		c.Gateway.EventFlowUrl = getenv("EVENTFLOW_URL", "http://eventflow:1880")
	}
	if c.Gateway.EmqxUrl == "" {
		c.Gateway.EmqxUrl = getenv("EMQX_URL", "http://emqx:18083")
	}
	if c.Storage.FileRoot == "" {
		c.Storage.FileRoot = getenv("FILESTORE_LOCAL_ROOT", "./data/files")
	}
	c.DataIngest.Enabled = getenvBool("DATAINGEST_ENABLED", c.DataIngest.Enabled)
	if c.DataIngest.MqttBrokers == "" {
		c.DataIngest.MqttBrokers = getenv("DATAINGEST_MQTT_BROKERS", "tcp://emqx:1883")
	}
	if c.DataIngest.MqttClientID == "" {
		c.DataIngest.MqttClientID = getenv("DATAINGEST_MQTT_CLIENT_ID", "edge-backend-sink")
	}
	if c.DataIngest.MqttUsername == "" {
		c.DataIngest.MqttUsername = getenv("DATAINGEST_MQTT_USERNAME", "")
	}
	if c.DataIngest.MqttPassword == "" {
		c.DataIngest.MqttPassword = getenv("DATAINGEST_MQTT_PASSWORD", "")
	}
	if c.DataIngest.MqttTopic == "" {
		c.DataIngest.MqttTopic = getenv("DATAINGEST_MQTT_TOPIC", "#")
	}
	if c.DataIngest.QueueSize <= 0 {
		c.DataIngest.QueueSize = getenvInt("DATAINGEST_QUEUE_SIZE", 20000)
	}
	if c.DataIngest.BatchSize <= 0 {
		c.DataIngest.BatchSize = getenvInt("DATAINGEST_BATCH_SIZE", 5000)
	}
	if c.DataIngest.FlushIntervalMs <= 0 {
		c.DataIngest.FlushIntervalMs = getenvInt("DATAINGEST_FLUSH_INTERVAL_MS", 1000)
	}
	c.Tier0Sdk.ApiHost = expandEnv(c.Tier0Sdk.ApiHost)
	c.Tier0Sdk.MqttHost = expandEnv(c.Tier0Sdk.MqttHost)
	c.Tier0Sdk.MqttPort = expandEnv(c.Tier0Sdk.MqttPort)
}
