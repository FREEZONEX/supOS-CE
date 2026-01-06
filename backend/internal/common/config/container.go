package config

// ContainerInfo represents container information
type ContainerInfo struct {
	Name        string         `json:"name" mapstructure:"name"`               // 容器名称
	Version     string         `json:"version" mapstructure:"version"`         // 容器版本
	Description string         `json:"description" mapstructure:"description"` // 容器描述
	EnvMap      map[string]any `json:"envMap" mapstructure:"env_map"`          // 环境变量
}

// ContainerEnv represents container environment variables
type ContainerEnv struct {
	ServiceIsShow      bool   `json:"serviceIsShow"`      // 服务是否显示
	ServiceLogo        string `json:"serviceLogo"`        // LOGO
	ServiceDescription string `json:"serviceDescription"` // 服务描述
	ServiceRedirectURL string `json:"serviceRedirectUrl"` // 高阶使用跳转路由
	ServiceAccount     string `json:"serviceAccount"`     // 帐号
	ServicePassword    string `json:"servicePassword"`    // 密码
}
