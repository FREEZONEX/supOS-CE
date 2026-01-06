package logic

import (
	"backend/internal/adapters/kong/listener"
	"backend/internal/adapters/kong/vo"
	"backend/internal/common/dto/protocol"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/go-resty/resty/v2"
	"github.com/magiconair/properties"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	routesPath   = "routes"
	servicesPath = "services"
	pluginsPath  = "plugins"
)

type (
	// KongLogic 封装了对 Kong Admin API 的直接调用和部分业务逻辑
	KongLogic struct {
		baseURL string
		client  *resty.Client
	}

	// InternalKongResponseVO Kong Admin API 的内部响应结构
	InternalKongResponseVO struct {
		Data []InternalKongVO `json:"data"`
	}

	// InternalKongVO Kong Admin API 的内部路由对象
	InternalKongVO struct {
		ID      string         `json:"id"`
		Name    string         `json:"name"`
		Paths   []string       `json:"paths"`
		Tags    []string       `json:"tags"`
		Service map[string]any `json:"service"`
	}

	// KongServiceRequest 创建 Kong Service 的请求体
	KongServiceRequest struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Enabled  bool   `json:"enabled"` // 服务是否启用
	}

	// KongRouteRequest 创建 Kong Route 的请求体
	KongRouteRequest struct {
		Name    string           `json:"name"`
		Paths   []string         `json:"paths"`
		Service KongRouteService `json:"service"`
		Tags    []string         `json:"tags,omitzero"`
	}

	KongRouteService struct {
		ID string `json:"id"`
	}

	// KongPluginRequest 创建 Kong Plugin 的请求体
	KongPluginRequest struct {
		Name      string         `json:"name"`
		Config    map[string]any `json:"config"`
		Enabled   bool           `json:"enabled"`
		Protocols []string       `json:"protocols"`
	}
)

var (
	kongLogic *KongLogic
	once      sync.Once
)

// GetKongLogic 创建 KongLogic 实例
func GetKongLogic(host string, port int) *KongLogic {
	once.Do(func() {
		kongLogic = &KongLogic{
			baseURL: fmt.Sprintf("http://%s:%d", host, port),
			client: resty.New().
				SetTimeout(30 * time.Second).
				SetRetryCount(3).
				SetRetryWaitTime(1 * time.Second),
		}
	})
	return kongLogic
}

// getRawRoutes 获取原始的路由列表
func (l *KongLogic) getRawRoutes(ctx context.Context) ([]InternalKongVO, error) {
	resp, err := l.client.R().
		SetResult(&InternalKongResponseVO{}).
		Get(l.baseURL + "/" + routesPath)

	if err != nil {
		logx.Errorf("request kong failed, error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		logx.Errorf("request kong failed, response: %s", resp.Body())
		return nil, fmt.Errorf("kong API error: %d", resp.StatusCode())
	}

	result := resp.Result().(*InternalKongResponseVO)
	return result.Data, nil
}

// QueryRoutes 查询并处理路由列表
func (l *KongLogic) QueryRoutes(ctx context.Context) ([]vo.RouteVO, error) {
	rawRoutes, err := l.getRawRoutes(ctx)
	if err != nil {
		return nil, err
	}

	// 获取本地菜单缓存，用于判断菜单是否被勾选
	localMenus := listener.GetLocalMenus()

	var routes []vo.RouteVO
	for _, kongVO := range rawRoutes {
		if l.IsMenu(kongVO.Tags) {
			// 基础菜单名称
			menuName := i18ns.LocalizeMsgWithCtx(ctx, kongVO.Name)
			showName := menuName

			// 解析标签并获取 showName
			tags, parsedShowName := l.parseTags(ctx, kongVO.Tags)
			if parsedShowName != "" {
				showName = parsedShowName
			}

			// 查询关联的服务信息
			service, _ := l.queryServiceById(kongVO.Service["id"].(string))

			rvo := vo.RouteVO{
				Name:     kongVO.Name,
				ShowName: showName,
				Tags:     tags,
				Service:  service,
			}

			if len(kongVO.Paths) > 0 {
				// 检查菜单是否在本地缓存中（被勾选）
				_, checked := localMenus[kongVO.Name]
				rvo.Menu = &vo.MenuVO{
					URL:    kongVO.Paths[0],
					Picked: checked,
				}
			}
			routes = append(routes, rvo)
		}
	}
	return routes, nil
}

// RouteList 查询简化的路由列表
func (l *KongLogic) RouteList(ctx context.Context) ([]vo.SimpleRouteVO, error) {
	rawRoutes, err := l.getRawRoutes(ctx)
	if err != nil {
		return nil, err
	}

	var routes []vo.SimpleRouteVO
	for _, route := range rawRoutes {
		for _, path := range route.Paths {
			routes = append(routes, vo.SimpleRouteVO{
				ID:   route.ID,
				Name: route.Name,
				URL:  path,
			})
		}
	}
	return routes, nil
}

// parseTags 解析 Kong 路由的标签
func (l *KongLogic) parseTags(ctx context.Context, tags []string) ([]*protocol.KeyValuePair[string], string) {
	var tagList []*protocol.KeyValuePair[string]
	var showName string
	for _, t := range tags {
		if tag, found := strings.CutPrefix(t, "parentName:"); found {
			tagShowName := i18ns.LocalizeMsgWithCtx(ctx, tag)
			tagList = append(tagList, &protocol.KeyValuePair[string]{Key: t, Value: tagShowName})
		} else if tag, found := strings.CutPrefix(t, "description:"); found {
			tagShowName := i18ns.LocalizeMsgWithCtx(ctx, tag)
			tagList = append(tagList, &protocol.KeyValuePair[string]{Key: t, Value: tagShowName})
		} else if tag, found := strings.CutPrefix(t, "homeParentName:"); found {
			tagShowName := i18ns.LocalizeMsgWithCtx(ctx, tag)
			tagList = append(tagList, &protocol.KeyValuePair[string]{Key: t, Value: tagShowName})
		} else if tag, found := strings.CutPrefix(t, "showName:"); found {
			showName = i18ns.LocalizeMsgWithCtx(ctx, tag)
			tagList = append(tagList, &protocol.KeyValuePair[string]{Key: t, Value: tag})
		} else {
			tagList = append(tagList, &protocol.KeyValuePair[string]{Key: "", Value: t})
		}
	}
	return tagList, showName
}

// queryServiceById 根据 ID 查询 Kong Service
func (l *KongLogic) queryServiceById(serviceId string) (*vo.ServiceResponseVO, error) {
	url := fmt.Sprintf("%s/%s/%s", l.baseURL, servicesPath, serviceId)
	result := &vo.ServiceResponseVO{}
	resp, err := l.client.R().
		SetResult(result).
		Get(url)

	if err != nil {
		logx.Errorf("request kong service failed, error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		logx.Errorf("request kong service failed, response: %s", resp.Body())
		return nil, nil
	}
	return result, nil
}

// QueryServiceJsonById 根据 name 查询 Kong Service
func (l *KongLogic) QueryServiceJsonById(serviceId string) (map[string]any, error) {
	url := fmt.Sprintf("%s/%s/%s", l.baseURL, servicesPath, serviceId)
	var result map[string]any
	resp, err := l.client.R().
		SetResult(&result).
		Get(url)
	if err != nil {
		logx.Errorf("request kong service failed, error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		logx.Errorf("request kong service failed, response: %s", resp.Body())
		return nil, nil
	}
	return result, nil
}

// IsMenu 检查路由是否为菜单
func (l *KongLogic) IsMenu(tags []string) bool {
	for _, tag := range tags {
		if tag == "menu" {
			return true
		}
	}
	return false
}

// CreateService 创建 Kong Service
func (l *KongLogic) CreateService(req *KongServiceRequest) (*vo.ServiceResponseVO, error) {
	url := l.baseURL + "/" + servicesPath
	result := &vo.ServiceResponseVO{}
	resp, err := l.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(result).
		Post(url)

	logx.Infof(">>>>>>>>>>>>kong create service URL： %s, params: %+v", url, req)
	if err != nil {
		logx.Errorf("request kong service failed, error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		logx.Errorf("request kong service failed, response: %s", resp.Body())
		return nil, fmt.Errorf("create service error, status: %d", resp.StatusCode())
	}
	return result, nil
}

// MarkMenu 标记菜单，将勾选的菜单持久化到本地文件
func (l *KongLogic) MarkMenu(routes []vo.MarkRouteRequestVO) error {
	// 1. 构建新的菜单映射
	newLocalMenus := make(map[string]string)
	for _, r := range routes {
		newLocalMenus[r.Name] = r.URL
	}

	// 2. 使用 properties 库构建并写入文件
	props := properties.LoadMap(newLocalMenus)

	// 3. 确保目录存在
	dir := filepath.Dir(listener.LocalMenuCheckedStoragePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logx.Errorf("failed to create directory %s: %v", dir, err)
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 4. 写入文件
	file, err := os.Create(listener.LocalMenuCheckedStoragePath)
	if err != nil {
		logx.Errorf("failed to create file %s: %v", listener.LocalMenuCheckedStoragePath, err)
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := props.Write(file, properties.UTF8); err != nil {
		logx.Errorf("failed to write properties to file: %v", err)
		return fmt.Errorf("failed to write properties: %w", err)
	}

	// 5. 更新内存缓存
	listener.UpdateLocalMenus(newLocalMenus)

	logx.Info("update menu cache success")
	return nil
}

// UpdateService 更新 Kong Service
func (l *KongLogic) UpdateService(id string, service map[string]any) (*vo.ServiceResponseVO, error) {
	url := fmt.Sprintf("%s/%s/%s", l.baseURL, servicesPath, id)
	result := &vo.ServiceResponseVO{}
	resp, err := l.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(service).
		SetResult(result).
		Put(url)

	if err != nil {
		logx.Errorf("request kong service failed, error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		logx.Errorf("request kong service failed, response: %s", resp.Body())
		return nil, fmt.Errorf("update service error")
	}
	return result, nil
}

// FetchRoute 查询单个路由，对应 public RoutResponseVO fetchRoute(String name)
func (l *KongLogic) FetchRoute(name string) (*vo.RoutResponseVO, error) {
	url := fmt.Sprintf("%s/%s/%s", l.baseURL, routesPath, name)
	result := &vo.RoutResponseVO{}
	resp, err := l.client.R().
		SetResult(result).
		Get(url)

	if err != nil {
		logx.Errorf("kong fetch route failed, error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		logx.Errorf("kong fetch route failed, response: %s", resp.Body())
		return nil, nil
	}
	return result, nil
}

// SearchRoute 按标签搜索路由，对应 public List<RoutResponseVO> searchRoute(List<String> tags)
func (l *KongLogic) SearchRoute(tags []string) ([]vo.RoutResponseVO, error) {
	url := fmt.Sprintf("%s/%s", l.baseURL, routesPath)
	var result InternalKongResponseVO
	resp, err := l.client.R().
		SetFormData(map[string]string{
			"tags": strings.Join(tags, ","),
		}).
		SetResult(&result).
		Get(url)

	if err != nil {
		logx.Errorf("kong search route failed, error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		logx.Errorf("kong search route failed, response: %s", resp.Body())
		return nil, nil
	}

	var routes []vo.RoutResponseVO
	for _, r := range result.Data {
		routes = append(routes, vo.RoutResponseVO{ID: r.ID, Name: r.Name})
	}
	return routes, nil
}

// CreateRoute 创建 Kong Route
func (l *KongLogic) CreateRoute(req *KongRouteRequest) (*vo.RoutResponseVO, error) {
	url := fmt.Sprintf("%s/%s", l.baseURL, routesPath)
	result := &vo.RoutResponseVO{}
	resp, err := l.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(result).
		Post(url)

	logx.Infof(">>>>>>>>>>>>kong create route URL： %s, params: %+v", url, req)
	if err != nil {
		logx.Errorf("kong create route failed, error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		logx.Errorf("kong create route failed, response: %s", resp.Body())
		return nil, fmt.Errorf("create route error, status: %d", resp.StatusCode())
	}
	return result, nil
}

// DeleteRoute 删除 Kong Route，对应 public void deleteRoute(String name)
func (l *KongLogic) DeleteRoute(name string) error {
	url := fmt.Sprintf("%s/%s/%s", l.baseURL, routesPath, name)
	logx.Infof(">>>>>>>>>>>>kong delete route URL： %s", url)
	resp, err := l.client.R().Delete(url)

	if err != nil {
		logx.Errorf("kong delete route failed, error: %v", err)
		return err
	}

	if resp.StatusCode() != http.StatusNoContent {
		logx.Errorf("kong delete route failed, response: %s", resp.Body())
		return fmt.Errorf("kong API error: %d", resp.StatusCode())
	}
	return nil
}

// UpdateRoute 更新 Kong Route
func (l *KongLogic) UpdateRoute(name string, req *KongRouteRequest) (*vo.RoutResponseVO, error) {
	url := fmt.Sprintf("%s/%s/%s", l.baseURL, routesPath, name)
	result := &vo.RoutResponseVO{}
	resp, err := l.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(result).
		Put(url)

	logx.Infof(">>>>>>>>>>>>kong update route URL： %s, params: %+v", url, req)
	if err != nil {
		logx.Errorf("kong update route failed, error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		logx.Errorf("kong update route failed, response: %s", resp.Body())
		return nil, fmt.Errorf("update route error, status: %d", resp.StatusCode())
	}
	return result, nil
}

// AddAPIKey 为 Route 添加 API Key 插件
func (l *KongLogic) AddAPIKey(name string) error {
	url := fmt.Sprintf("%s/%s/%s/%s", l.baseURL, routesPath, name, pluginsPath)
	pluginReq := &KongPluginRequest{
		Name: "key-auth",
		Config: map[string]any{
			"key_names":        []string{"apikey"},
			"key_in_header":    true,
			"key_in_query":     true,
			"anonymous":        nil,
			"run_on_preflight": true,
			"hide_credentials": false,
			"key_in_body":      false,
			"realm":            nil,
		},
		Enabled:   true,
		Protocols: []string{"grpc", "grpcs", "http", "https"},
	}

	resp, err := l.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(pluginReq).
		Post(url)

	logx.Infof(">>>>>>>>>>>>kong addApiKey URL： %s, params: %+v", url, pluginReq)
	if err != nil {
		logx.Errorf("kong addApiKey failed, error: %v", err)
		return err
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		logx.Errorf("kong addApiKey failed, response: %s", resp.Body())
		return fmt.Errorf("add api key error, status: %d", resp.StatusCode())
	}
	return nil
}
