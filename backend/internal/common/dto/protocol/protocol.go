package protocol

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// KeyValuePair 键值对，对应 Java 的泛型类 KeyValuePair<T>
type KeyValuePair[T any] struct {
	Key   string `json:"key" validate:"required"`
	Value T      `json:"value"`
}

// NewKeyValuePair 创建 KeyValuePair
func NewKeyValuePair[T any](key string, value T) *KeyValuePair[T] {
	return &KeyValuePair[T]{Key: key, Value: value}
}

// BaseConfigDTO 基础协议配置
type BaseConfigDTO struct {
	Protocol string `json:"protocol"`
}

// BaseServerConfigDTO 基础服务器配置
type BaseServerConfigDTO struct {
	Host string `json:"host" validate:"required"`
	Port string `json:"port"` // Java 中是 String 类型
}

// SimpleModelDTO 简单模型
type SimpleModelDTO struct {
	Alias string `json:"alias"`
	Topic string `json:"topic"`
}

// MappingDTO 映射配置
type MappingDTO struct {
	Alias string `json:"alias"`
	Name  string `json:"name"` // 属性名称
}

// RateDTO 速率配置
type RateDTO struct {
	Value int64  `json:"value" validate:"min=1"`
	Unit  string `json:"unit,omitzero"` // ms, s, m, h
}

// GetUnit 获取单位，默认为 "s"
func (r *RateDTO) GetUnit() string {
	if r.Unit != "" {
		return r.Unit
	}
	return "s"
}

// GetSeconds 获取秒数
func (r *RateDTO) GetSeconds() int64 {
	if r.Unit == "" {
		return 60 // 1分钟
	}
	switch r.Unit {
	case "ms":
		v := r.Value / 1000
		if v <= 0 {
			return 1
		}
		return v
	case "s":
		return r.Value
	case "m":
		return r.Value * 60
	case "h":
		return r.Value * 3600
	default:
		return 60
	}
}

// OpcUAConfigDTO OPC UA 配置
type OpcUAConfigDTO struct {
	BaseConfigDTO
	ServerName string                `json:"serverName,omitzero"`
	PollRate   *RateDTO              `json:"pollRate" validate:"required"`
	Server     *OpcuaServerConfigDTO `json:"server,omitzero"`
}

// GetServerName 获取服务器名称
func (o *OpcUAConfigDTO) GetServerName() string {
	if o.ServerName != "" {
		return o.ServerName
	}
	if o.Server != nil {
		return o.Server.GetEndpoint()
	}
	return ""
}

// OpcuaServerConfigDTO OPC UA 服务器配置
type OpcuaServerConfigDTO struct {
	BaseServerConfigDTO
	Location string `json:"location,omitzero"`
}

// GetEndpoint 获取 OPC UA 端点地址
func (o *OpcuaServerConfigDTO) GetEndpoint() string {
	lt := ""
	if o.Location != "" {
		if strings.HasPrefix(o.Location, "/") {
			lt = o.Location
		} else {
			lt = "/" + o.Location
		}
	}
	endpoint := fmt.Sprintf("opc.tcp://%s:%s%s", o.Host, o.Port, lt)
	return strings.TrimSpace(endpoint)
}

// OpcDAConfigDTO OPC DA 配置
type OpcDAConfigDTO struct {
	BaseConfigDTO
	ServerName string                `json:"serverName,omitzero"`
	PollRate   *RateDTO              `json:"pollRate" validate:"required"`
	Server     *OpcdaServerConfigDTO `json:"server,omitzero"`
}

// GetServerName 获取服务器名称
func (o *OpcDAConfigDTO) GetServerName() string {
	if o.ServerName != "" {
		return o.ServerName
	}
	if o.Server != nil {
		return o.Server.Host
	}
	return ""
}

// OpcdaServerConfigDTO OPC DA 服务器配置
type OpcdaServerConfigDTO struct {
	BaseServerConfigDTO
	Domain   string `json:"domain,omitzero"`
	Account  string `json:"account,omitzero"`
	Password string `json:"password,omitzero"`
	Clsid    string `json:"clsid,omitzero"`
	Timeout  int64  `json:"timeout,omitzero"` // ms
}

// GetTimeout 获取超时时间，默认 5000ms
func (o *OpcdaServerConfigDTO) GetTimeout() int64 {
	if o.Timeout > 0 {
		return o.Timeout
	}
	return 5000
}

// MqttConfigDTO MQTT 配置
type MqttConfigDTO struct {
	BaseConfigDTO
	Server     *MqttServerConfigDTO `json:"server,omitzero"`
	ServerName string               `json:"serverName,omitzero"`
	InputName  string               `json:"inputName,omitzero"`
	InputTopic string               `json:"inputTopic,omitzero"`
}

// GetServerName 获取服务器名称
func (m *MqttConfigDTO) GetServerName() string {
	if m.ServerName != "" {
		return m.ServerName
	}
	if m.Server != nil {
		return m.Server.Host + ":" + m.Server.Port
	}
	return ""
}

// GetInputName 获取输入名称
func (m *MqttConfigDTO) GetInputName() string {
	return m.InputName
}

// MqttServerConfigDTO MQTT 服务器配置
type MqttServerConfigDTO struct {
	BaseServerConfigDTO
	Username string `json:"username,omitzero"`
	Password string `json:"password,omitzero"`
}

// ModbusConfigDTO Modbus 配置
type ModbusConfigDTO struct {
	BaseConfigDTO
	ServerName string                 `json:"serverName,omitzero"`
	UnitID     string                 `json:"unitId" validate:"required"`
	FC         string                 `json:"fc" validate:"required"` // FC1->Coil, FC2->Input, FC3->HoldingRegister, FC4->InputRegister
	Address    string                 `json:"address" validate:"required"`
	Quantity   string                 `json:"quantity,omitzero"`
	PollRate   *RateDTO               `json:"pollRate" validate:"required"`
	Server     *ModbusServerConfigDTO `json:"server" validate:"required"`
}

// GetServerName 获取服务器名称
func (m *ModbusConfigDTO) GetServerName() string {
	if m.ServerName != "" {
		return m.ServerName
	}
	if m.Server == nil {
		return ""
	}
	return m.Server.Host + ":" + m.Server.Port + "@" + m.UnitID
}

// ModbusServerConfigDTO Modbus 服务器配置
type ModbusServerConfigDTO struct {
	BaseServerConfigDTO
	UnitID string `json:"unitId,omitzero"`
}

// RestConfigDTO REST 配置
type RestConfigDTO struct {
	BaseConfigDTO
	ServerName string                  `json:"serverName,omitzero"`
	SyncRate   *RateDTO                `json:"syncRate" validate:"required"` // 同步第三方接口的频率
	Params     []*KeyValuePair[string] `json:"params,omitzero"`              // 查询参数
	Headers    json.RawMessage         `json:"headers,omitzero"`             // JSONArray
	PageDef    *PageDef                `json:"pageDef,omitzero"`             // 分页参数
	Path       string                  `json:"path,omitzero"`
	HTTPS      bool                    `json:"https"`
	Body       string                  `json:"body,omitzero"`
	Method     string                  `json:"method" validate:"required"`
	Server     *RestServerConfigDTO    `json:"server,omitzero"`
	FullURL    string                  `json:"fullUrl,omitzero"` // 完整地址
}

// GetBody 获取请求体，默认 "{}"
func (r *RestConfigDTO) GetBody() string {
	if r.Body != "" {
		return r.Body
	}
	return "{}"
}

// GetURL 获取 URL
func (r *RestConfigDTO) GetURL() string {
	if r.FullURL != "" {
		return r.FullURL
	}
	if r.Server != nil {
		proto := "http"
		if r.HTTPS {
			proto = "https"
		}
		urlStr := proto + "://" + r.Server.Host
		if r.Server.Port != "" {
			urlStr += ":" + r.Server.Port
		}
		if r.Path != "" {
			if strings.HasPrefix(r.Path, "/") {
				urlStr += r.Path
			} else {
				urlStr += "/" + r.Path
			}
		}
		return urlStr
	}
	return ""
}

// SetFullURL 设置完整 URL 并解析
func (r *RestConfigDTO) SetFullURL(fullURL string) error {
	r.FullURL = fullURL
	if fullURL == "" {
		return nil
	}

	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return err
	}

	r.Server = &RestServerConfigDTO{}
	r.Server.Host = parsedURL.Host
	if parsedURL.Port() != "" {
		r.Server.Port = parsedURL.Port()
	}
	r.HTTPS = parsedURL.Scheme == "https"
	r.Path = parsedURL.Path

	// 解析查询参数
	if parsedURL.RawQuery != "" {
		pageK, offsetK := "", ""
		if r.PageDef != nil {
			if r.PageDef.Start != nil {
				pageK = r.PageDef.Start.Key
			}
			if r.PageDef.Offset != nil {
				offsetK = r.PageDef.Offset.Key
			}
		}

		queryParams := parsedURL.Query()
		for key, values := range queryParams {
			if key != pageK && key != offsetK && len(values) > 0 {
				r.Params = append(r.Params, NewKeyValuePair(key, values[0]))
			}
		}
	}

	return nil
}

// SetPageDef 设置分页参数
func (r *RestConfigDTO) SetPageDef(pageDef *PageDef) {
	r.PageDef = pageDef
	if pageDef == nil || len(r.Params) == 0 {
		return
	}

	pageK, offsetK := "", ""
	if pageDef.Start != nil && pageDef.Start.Key != "" {
		pageK = pageDef.Start.Key
	}
	if pageDef.Offset != nil && pageDef.Offset.Key != "" {
		offsetK = pageDef.Offset.Key
	}

	// 从 params 中移除分页参数
	newParams := make([]*KeyValuePair[string], 0, len(r.Params))
	for _, pair := range r.Params {
		switch pair.Key {
		case pageK:
			if pageDef.Start != nil && pageDef.Start.Value == "" {
				pageDef.Start.Value = pair.Value
			}
		case offsetK:
			if pageDef.Offset != nil && pageDef.Offset.Value == "" {
				pageDef.Offset.Value = pair.Value
			}
		default:
			newParams = append(newParams, pair)
		}
	}
	r.Params = newParams
}

// GainFullURL 获取完整 URL（包含查询参数）
func (r *RestConfigDTO) GainFullURL() string {
	if r.FullURL != "" {
		return r.FullURL
	}

	baseURL := r.GetURL()
	if baseURL == "" || len(r.Params) == 0 {
		return baseURL
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	query := parsedURL.Query()
	for _, param := range r.Params {
		// 将 any 类型的 Value 转换为 string
		valueStr := ""
		if param.Value != "" {
			valueStr = fmt.Sprintf("%v", param.Value)
		}
		query.Add(param.Key, valueStr)
	}
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}

// PageDef 分页参数定义
type PageDef struct {
	Start  *KeyValuePair[string] `json:"start" validate:"required"`
	Offset *KeyValuePair[string] `json:"offset,omitzero"`
}

// RestServerConfigDTO REST 服务器配置
type RestServerConfigDTO struct {
	BaseServerConfigDTO
}

// ICMPConfigDTO ICMP 配置
type ICMPConfigDTO struct {
	BaseConfigDTO
	Interval int                  `json:"interval" validate:"min=1"` // ping时间间隔
	Timeout  int                  `json:"timeout" validate:"min=1"`  // 超时时间，单位秒
	Retry    int                  `json:"retry" validate:"min=1"`    // 重试次数
	Server   *BaseServerConfigDTO `json:"server,omitzero"`
}

// ProtocolTagEnums 位号枚举
type ProtocolTagEnums struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

// Equal 检查两个 ProtocolTagEnums 实例是否逻辑相等。
func (p *ProtocolTagEnums) Equal(other *ProtocolTagEnums) bool {
	// 如果两个指针都指向同一个实例，或者都为 nil，则它们相等。
	if p == other {
		return true
	}
	// 如果其中一个为 nil，则它们不相等。
	if p == nil || other == nil {
		return false
	}
	// 比较核心字段的值。
	return p.Name == other.Name && p.DataType == other.DataType
}
