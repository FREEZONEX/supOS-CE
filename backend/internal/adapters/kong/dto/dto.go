package dto

import "mime/multipart"

type MenuDto struct {
	// 所属服务名称
	ServiceName string `json:"serviceName,omitzero" form:"serviceName"`
	// 菜单名称，唯一标识
	Name string `json:"name" form:"name"`
	// 显示名称，支持国际化
	ShowName string `json:"showName" form:"showName"`
	// 描述
	Description string `json:"description,omitzero" form:"description"`
	// 基础 URL
	BaseURL string `json:"baseUrl" form:"baseUrl"`
	// 打开方式 (0: iframe, 1: 新页面)
	OpenType int `json:"openType" form:"openType"`
	// 菜单图标文件
	Icon *multipart.FileHeader `json:"-" form:"icon"`
	// 是否是菜单
	IsMenu bool `json:"isMenu,omitzero"`
	// 标签
	Tags []string `json:"tags,omitzero"`
}

type ResourceQuery struct {
	// 父级ID
	ParentID int64 `json:"parentId,omitzero" form:"parentId"`
	// 资源类型 (1-目录 2-菜单 3-按钮 4-Tab 5-子菜单)
	Type int `json:"type,omitzero" form:"type"`
}
