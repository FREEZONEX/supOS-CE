package types

import "backend/share/app/model"

// InstalledAppsResponse 已安装应用列表响应
type InstalledAppsResponse struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg,optional"`
	Data []model.NewFeatureModel `json:"data"`
}

// AppDetailResponse 应用详情响应
type AppDetailResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message,optional"`
	Data    AppDetail `json:"data"`
}

// InstallAppResponse 安装应用响应
type InstallAppResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message,optional"`
	Data    InstallResult `json:"data"`
}

// UninstallAppResponse 卸载应用响应
type UninstallAppResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message,optional"`
	Data    UninstallResult `json:"data"`
}

// UpdateAppResponse 更新应用响应
type UpdateAppResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message,optional"`
	Data    UpdateResult `json:"data"`
}

// AppInfo 应用信息
type AppInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageSource string `json:"imageSource"`
	MenuUrl     string `json:"menuUrl"`
	ComposeYaml string `json:"composeYaml"`
	InstallTime string `json:"installTime,optional"`
	Status      string `json:"status,optional"`
}

// AppDetail 应用详情
type AppDetail struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageSource string `json:"imageSource"`
	ImagePath   string `json:"imagePath,omitempty"`
	ImageUrl    string `json:"imageUrl,omitempty"`
	IconPath    string `json:"iconPath,omitempty"`
	MenuUrl     string `json:"menuUrl"`
	ComposeYaml string `json:"composeYaml"`
	RouterTrim  bool   `json:"routerTrim"`
	KeepHost    bool   `json:"keepHost"`
	InstallTime string `json:"installTime,omitempty"`
	Status      string `json:"status,optional"`
}

// InstallResult 安装结果
type InstallResult struct {
	Name        string `json:"name"`
	Success     bool   `json:"success"`
	Message     string `json:"message,optional"`
	InstallTime string `json:"installTime,optional"`
}

// UninstallResult 卸载结果
type UninstallResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message,optional"`
}

// UpdateResult 更新结果
type UpdateResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message,optional"`
}

// GetAppRequest 获取应用请求
type GetAppRequest struct {
	Name string `path:"name"`
}

// SearchAppsRequest 搜索应用请求
type SearchAppsRequest struct {
	Keyword string `form:"keyword,optional"`
}

// InstallAppRequest 安装应用请求
type InstallAppRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description,optional"`
	ImagePath   string           `json:"imagePath,optional"`
	ImageUrl    string           `json:"imageUrl,optional"`
	Menu        *model.MenuModel `json:"menu"`
	ComposeYaml string           `json:"composeYaml"`
}

// UninstallAppRequest 卸载应用请求
type UninstallAppRequest struct {
	Name string `path:"name"`
}
