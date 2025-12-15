package config

// SystemConfig represents system-level configuration
type SystemConfig struct {
	AppTitle     string `mapstructure:"app_title"`     // 应用标题
	Version      string `mapstructure:"version"`       // 系统版本
	Lang         string `mapstructure:"lang"`          // 语言 (en-US, zh-CN)
	AuthEnable   bool   `mapstructure:"auth_enable"`   // 是否开启keycloak校验
	LLMType      string `mapstructure:"llm_type"`      // 大语言模型类型
	PlatformType string `mapstructure:"platform_type"` // 基础平台类型
	EntranceURL  string `mapstructure:"entrance_url"`  // 系统入口地址
	LoginPath    string `mapstructure:"login_path"`    // 登录页url

	// MQTT configuration
	MQTTTCPPort          int  `mapstructure:"mqtt_tcp_port"`           // MQTT TCP端口
	MQTTWebsocketTSLPort int  `mapstructure:"mqtt_websocket_tsl_port"` // MQTT WebSocket加密端口
	MultipleTopic        bool `mapstructure:"multiple_topic"`          // 单双topic
	UseAliasPathAsTopic  bool `mapstructure:"use_alias_path_as_topic"` // 是否使用别名作为MQTT topic

	// Field names
	QualityName   string `mapstructure:"quality_name"`   // 质量码字段名称
	TimestampName string `mapstructure:"timestamp_name"` // 系统时间字段名称

	// Features
	LazyTree                 bool `mapstructure:"lazy_tree"`                  // 是否启用UNS树懒加载模式
	LDAPEnable               bool `mapstructure:"ldap_enable"`                // 是否启用LDAP用户体系
	EnableAutoCategorization bool `mapstructure:"enable_auto_categorization"` //是否开启文件自动归类
	// Container configuration
	ContainerMap map[string]*ContainerInfo `mapstructure:"containers"` // 系统容器

}

// NewSystemConfig returns a SystemConfig with default values
func NewSystemConfig() *SystemConfig {
	return &SystemConfig{
		Lang:                     "zh-CN",
		AuthEnable:               false,
		LLMType:                  "ollama",
		MQTTTCPPort:              1883,
		MQTTWebsocketTSLPort:     8084,
		LoginPath:                "/tier0-login",
		PlatformType:             "linux",
		QualityName:              "status",
		TimestampName:            "timeStamp",
		LazyTree:                 false,
		LDAPEnable:               false,
		EnableAutoCategorization: true,
		ContainerMap:             make(map[string]*ContainerInfo),
	}
}
