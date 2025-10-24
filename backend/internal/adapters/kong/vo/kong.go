package vo

import "backend/internal/common/dto/protocol"

type RouteVO struct {
	// 名称，可能是国际化 key
	Name string `json:"name"`
	// 显示名称
	ShowName string `json:"showName"`
	// 菜单信息
	Menu *MenuVO `json:"menu,omitzero"`
	// 标签
	Tags []*protocol.KeyValuePair[string] `json:"tags,omitzero"`
	// 关联的服务
	Service *ServiceResponseVO `json:"service,omitzero"`
}

type MenuVO struct {
	// 菜单 URL
	URL string `json:"url"`
	// 是否被选中
	Picked bool `json:"picked,omitzero"`
}

type SimpleRouteVO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ServiceResponseVO struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Host     string   `json:"host"`
	Path     string   `json:"path"`
	Port     int      `json:"port"`
	Protocol string   `json:"protocol"`
	Tags     []string `json:"tags"`
}

// RoutResponseVO Kong Admin API 对 Route 的响应体
type RoutResponseVO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MarkRouteRequestVO 用于标记菜单的请求体
type MarkRouteRequestVO struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ResourceVO 用于展示资源信息
type ResourceVO struct {
	// 主键ID
	ID string `json:"id,omitzero"`
	// 父级ID
	ParentID string `json:"parentId,omitzero"`
	// 资源类型 (1-目录 2-菜单 3-按钮 4-Tab 5-子菜单)
	Type int `json:"type,omitzero"`
	// 资源编码
	Code string `json:"code,omitzero"`
	// 名称国际化 code
	NameCode string `json:"nameCode,omitzero"`
	// 显示名称
	ShowName string `json:"showName,omitzero"`
	// 路由来源 (1-手工 2-Kong)
	RouteSource int `json:"routeSource,omitzero"`
	// 地址
	URL string `json:"url,omitzero"`
	// 类型 (1-内部地址 2-外部链接)
	URLType int `json:"urlType,omitzero"`
	// 打开方式 (0-当前页面跳转 1-新窗口打开)
	OpenType int `json:"openType,omitzero"`
	// 图标
	Icon string `json:"icon,omitzero"`
	// 描述国际化 Key
	DescriptionCode string `json:"descriptionCode,omitzero"`
	// 描述内容
	ShowDescription string `json:"showDescription,omitzero"`
	// 排序
	Sort int `json:"sort,omitzero"`
	// 是否可编辑
	EditEnable bool `json:"editEnable,omitzero"`
	// 是否显示在首页
	HomeEnable bool `json:"homeEnable,omitzero"`
	// 是否固定
	Fixed bool `json:"fixed,omitzero"`
	// 启用状态
	Enable bool `json:"enable,omitzero"`
	// 更新时间
	UpdateAt string `json:"updateAt,omitzero"`
	// 创建时间
	CreateAt string `json:"createAt,omitzero"`
}

// ResultVO 通用的 API 响应结构
type ResultVO[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitzero"`
	Data T      `json:"data,omitzero"`
}

// Success 成功响应（带数据）
func Success[T any](data T) *ResultVO[T] {
	return &ResultVO[T]{Code: 200, Data: data}
}

// SuccessWithMsg 成功响应（带消息）
func SuccessWithMsg(msg string) *ResultVO[any] {
	return &ResultVO[any]{Code: 200, Msg: msg}
}

// Fail 失败响应
func Fail(msg string) *ResultVO[any] {
	return &ResultVO[any]{Code: 500, Msg: msg}
}
