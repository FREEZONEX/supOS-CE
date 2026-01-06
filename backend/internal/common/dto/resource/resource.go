package resource

// SaveResourceDto 保存资源 DTO
type SaveResourceDto struct {
	ID          int64             `json:"id,omitzero"`
	Code        string            `json:"code" validate:"required"`
	Name        string            `json:"name" validate:"required"`
	ParentID    int64             `json:"parentId,omitzero"`
	Type        int               `json:"type" validate:"required"` // 1-目录 2-菜单 3-按钮 4-Tab
	Source      string            `json:"source,omitzero"`
	RouteSource int               `json:"routeSource,omitzero"` // 1-手工 2-Kong
	URL         string            `json:"url,omitzero"`
	URLType     int               `json:"urlType,omitzero"` // 1-内部地址 2-外部链接
	OpenType    int               `json:"openType"`         // 0-当前页面跳转 1-新窗口打开
	Icon        string            `json:"icon,omitzero"`
	Description string            `json:"description,omitzero"`
	Sort        int               `json:"sort,omitzero"`
	EditEnable  bool              `json:"editEnable,omitzero"`
	HomeEnable  bool              `json:"homeEnable,omitzero"`
	Fixed       bool              `json:"fixed,omitzero"`
	Enable      bool              `json:"enable,omitzero"`
	Children    []SaveResourceDto `json:"children,omitzero"`
}

// SaveResource4ExternalDto 外部保存资源 DTO（带国际化支持）
type SaveResource4ExternalDto struct {
	Code            string `json:"code" validate:"required"`
	NameCode        string `json:"nameCode" validate:"required"`
	ParentID        int64  `json:"parentId,omitzero"`
	Type            int    `json:"type" validate:"required"` // 1-目录 2-菜单 3-按钮 4-Tab
	Source          string `json:"source" validate:"required"`
	URL             string `json:"url,omitzero"`
	URLType         int    `json:"urlType,omitzero"` // 1-内部地址 2-外部链接
	OpenType        int    `json:"openType"`         // 0-当前页面跳转 1-新窗口打开
	Icon            string `json:"icon,omitzero"`
	DescriptionCode string `json:"descriptionCode,omitzero"`
	Sort            int    `json:"sort,omitzero"`
	EditEnable      bool   `json:"editEnable,omitzero"`
	HomeEnable      bool   `json:"homeEnable,omitzero"`
	Fixed           bool   `json:"fixed,omitzero"`
	Enable          bool   `json:"enable,omitzero"`
}

// BatchUpdateResourceDto 批量更新资源 DTO
type BatchUpdateResourceDto struct {
	ID              int64  `json:"id,omitzero"`
	Code            string `json:"code" validate:"required"`
	Name            string `json:"name" validate:"required"`
	NameCode        string `json:"nameCode,omitzero"`
	DescriptionCode string `json:"descriptionCode,omitzero"`
	ParentID        int64  `json:"parentId,omitzero"`
	Type            int    `json:"type" validate:"required"` // 1-目录 2-菜单 3-按钮 4-Tab
	Source          string `json:"source,omitzero"`
	RouteSource     int    `json:"routeSource,omitzero"` // 1-手工 2-Kong
	URL             string `json:"url,omitzero"`
	URLType         int    `json:"urlType,omitzero"` // 1-内部地址 2-外部链接
	OpenType        int    `json:"openType"`         // 0-当前页面跳转 1-新窗口打开
	Icon            string `json:"icon,omitzero"`
	Description     string `json:"description,omitzero"`
	Sort            int    `json:"sort,omitzero"`
	EditEnable      bool   `json:"editEnable,omitzero"`
	HomeEnable      bool   `json:"homeEnable,omitzero"`
	Fixed           bool   `json:"fixed,omitzero"`
	Enable          bool   `json:"enable,omitzero"`
}
