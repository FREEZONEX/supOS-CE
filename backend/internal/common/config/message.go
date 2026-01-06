package config

// MessageSourceConfig represents message source configuration for i18n.
type MessageSourceConfig struct {
	Basename       string `mapstructure:"basename"`         // 资源文件基础名称（可以是逗号分隔的多个）
	Encoding       string `mapstructure:"encoding"`         // 编码格式
	DefaultLocale  string `mapstructure:"default_locale"`   // 默认语言环境
	IsFileBasename bool   `mapstructure:"is_file_basename"` // 是否是 file: 开头的外置资源文件
}

// NewMessageSourceConfig returns default message source configuration
func NewMessageSourceConfig() *MessageSourceConfig {
	return &MessageSourceConfig{
		Basename:       "messages",
		Encoding:       "UTF-8",
		DefaultLocale:  "en-US", // Java uses Locale.US
		IsFileBasename: false,
	}
}

// GetBasenames returns basenames as a slice (split by comma if multiple)
func (m *MessageSourceConfig) GetBasenames() []string {
	if m.Basename == "" {
		return []string{"messages"}
	}
	// Split the basename string by commas.
	return splitByComma(m.Basename)
}

// IsReloadable checks if this should use a reloadable message source.
func (m *MessageSourceConfig) IsReloadable() bool {
	return m.IsFileBasename || len(m.Basename) > 5 && m.Basename[:5] == "file:"
}

func splitByComma(s string) []string {
	result := []string{}
	current := ""
	for _, c := range s {
		if c == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
