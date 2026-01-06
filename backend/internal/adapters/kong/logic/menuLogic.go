package logic

import (
	"backend/internal/adapters/kong/dto"
	"backend/internal/common/constants"
	"backend/internal/common/errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

var iconRootPath = filepath.Join(constants.RootPath, "system", "resource", "supos")

// MenuLogic 封装了菜单和路由相关的核心业务逻辑
type MenuLogic struct {
	kongLogic *KongLogic
}

// NewMenuLogic 创建 MenuLogic 实例
func NewMenuLogic(host string, port int) *MenuLogic {
	return &MenuLogic{
		kongLogic: GetKongLogic(host, port),
	}
}

// CreateRouteWithNoService 创建路由（不带服务）
func (l *MenuLogic) CreateRouteWithNoService(menuDto *dto.MenuDto, needApiKey, updateService bool) error {
	// 1. 获取或创建服务（不提供 protocol, host, port 信息）
	serviceID, _, err := l.fetchOrCreateService(menuDto.ServiceName, updateService, "", "", 0)
	if err != nil {
		return err
	}
	// 2. 创建或更新路由
	return l.createOrUpdateRoute(serviceID, menuDto.Name, menuDto.BaseURL, menuDto.Tags, needApiKey)
}

// createOrUpdateRoute 是创建或更新路由的通用逻辑
func (l *MenuLogic) createOrUpdateRoute(serviceID, routeName, path string, tags []string, needApiKey bool) error {
	existRoute, err := l.kongLogic.FetchRoute(routeName)
	if err != nil {
		return errors.NewBuzError(500, "menu.save.failed")
	}

	routeReq := &KongRouteRequest{
		Name:    routeName,
		Paths:   []string{path},
		Service: KongRouteService{ID: serviceID},
		Tags:    tags,
	}

	if existRoute != nil {
		_, err = l.kongLogic.UpdateRoute(routeName, routeReq)
	} else {
		_, err = l.kongLogic.CreateRoute(routeReq)
		if err == nil && needApiKey {
			err = l.kongLogic.AddAPIKey(routeName)
		}
	}

	if err != nil {
		logx.Errorf("failed to save route: %v", err)
		return errors.NewBuzError(500, "menu.save.failed")
	}
	return nil
}

// CreateRoute 创建路由（带服务），对应 public void createRoute(...)
func (l *MenuLogic) CreateRoute(menuDto *dto.MenuDto, needApiKey, updateService bool) error {
	// 1. 解析和校验 BaseUrl
	parsedURL, err := url.Parse(menuDto.BaseURL)
	if err != nil {
		return errors.NewBuzError(500, "menu.baseurl.invalid")
	}
	host := parsedURL.Hostname()
	protocol := parsedURL.Scheme
	port := parsedURL.Port()
	if port == "" {
		switch protocol {
		case "https":
			port = "443"
		case "http":
			port = "80"
		}
	}
	path := parsedURL.Path
	if path == "" || path == "/" {
		if parsedURL.RawQuery != "" {
			path = "/?" + parsedURL.RawQuery
		} else {
			path = "/"
		}
	} else if parsedURL.RawQuery != "" {
		path = path + "?" + parsedURL.RawQuery
	}
	// 校验 host
	if _, err := net.LookupHost(host); err != nil {
		return errors.NewBuzError(500, "menu creation error: UnknownHost")
	}

	// 2. 获取或创建服务
	portInt := cast.ToInt(port)
	serviceID, _, err := l.fetchOrCreateService(menuDto.ServiceName, updateService, protocol, host, portInt)
	if err != nil {
		return err
	}

	// 3. 创建或更新路由
	return l.createOrUpdateRoute(serviceID, menuDto.Name, path, menuDto.Tags, needApiKey)
}

// CreateMenu 创建菜单
func (l *MenuLogic) CreateMenu(menuDto *dto.MenuDto, updateService bool) error {
	// 1. 解析和校验 BaseUrl
	parsedURL, err := url.Parse(menuDto.BaseURL)
	if err != nil {
		return errors.NewBuzError(500, "menu.baseurl.invalid")
	}
	host := parsedURL.Hostname()
	protocol := parsedURL.Scheme
	portStr := parsedURL.Port()
	if portStr == "" {
		switch protocol {
		case "https":
			portStr = "443"
		case "http":
			portStr = "80"
		}
	}
	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	// 校验 host 是否可达
	if _, err := net.LookupHost(host); err != nil {
		logx.Errorf("menu creation error, unknown host: %v", err)
		return errors.NewBuzError(500, "menu creation error: UnknownHost")
	}

	// 2. 获取或创建服务
	portInt := cast.ToInt(portStr)
	serviceID, serviceName, err := l.fetchOrCreateService(menuDto.ServiceName, updateService, protocol, host, portInt)
	if err != nil {
		return err
	}

	// 3. 处理图标文件
	var iconName string
	if menuDto.Icon != nil {
		ext := filepath.Ext(menuDto.Icon.Filename)
		iconName = menuDto.Name + ext

		if err := l.saveIconFile(menuDto.Icon, iconName); err != nil {
			return err
		}
	}

	// 4. 封装标签
	tags := l.buildMenuTags(menuDto, iconName)

	// 5. 检查显示名是否已存在
	existShowNames, err := l.kongLogic.SearchRoute([]string{fmt.Sprintf("showName:%s", menuDto.ShowName), "menu"})
	if err != nil {
		return errors.NewBuzError(500, "menu.save.failed")
	}
	for _, existShowName := range existShowNames {
		if existShowName.Name != menuDto.Name {
			return errors.NewBuzError(500, "menu.showname.exist")
		}
	}

	// 6. 根据 openType 调整路径
	if menuDto.OpenType != 2 {
		path = "/third-apps/" + serviceName + path
	}

	// 7. 创建或更新路由
	existRoute, err := l.kongLogic.FetchRoute(menuDto.Name)
	if err != nil {
		return errors.NewBuzError(500, "menu.save.failed")
	}

	routeReq := &KongRouteRequest{
		Name:    menuDto.Name,
		Paths:   []string{path},
		Service: KongRouteService{ID: serviceID},
		Tags:    tags,
	}

	if existRoute != nil {
		// 更新
		_, err = l.kongLogic.UpdateRoute(menuDto.Name, routeReq)
	} else {
		// 创建
		_, err = l.kongLogic.CreateRoute(routeReq)
		// 注意：Java 代码中 app 变量被注释掉了（第203行），所以 addApiKey 永远不会执行
		// 因此这里也不应该调用 AddAPIKey
		// if err == nil && menuDto.OpenType == 2 { // ❌ Java 中这段逻辑是死代码
		// 	err = l.kongLogic.AddAPIKey(menuDto.Name)
		// }
	}

	if err != nil {
		logx.Errorf("failed to save route: %v", err)
		return errors.NewBuzError(500, "menu.save.failed")
	}

	return nil
}

// DeleteMenu 删除菜单
func (l *MenuLogic) DeleteMenu(name string) error {
	existRoute, err := l.kongLogic.FetchRoute(name)
	if err != nil {
		logx.Errorf("failed to check route %s before deletion: %v", name, err)
	}
	if existRoute != nil {
		return l.kongLogic.DeleteRoute(name)
	}
	return nil
}

// private helper functions

func (l *MenuLogic) saveIconFile(fileHeader *multipart.FileHeader, iconName string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return errors.NewBuzError(500, "menu.icon.save.failed")
	}
	defer src.Close()

	// 确保目录存在
	if err := os.MkdirAll(iconRootPath, 0755); err != nil {
		logx.Errorf("failed to create icon directory: %v", err)
		return errors.NewBuzError(500, "menu.icon.save.failed")
	}

	dst, err := os.Create(filepath.Join(iconRootPath, iconName))
	if err != nil {
		logx.Errorf("failed to create icon file: %v", err)
		return errors.NewBuzError(500, "menu.icon.save.failed")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		logx.Errorf("failed to write icon file: %v", err)
		return errors.NewBuzError(500, "menu.icon.save.failed")
	}
	return nil
}

func (l *MenuLogic) buildMenuTags(menuDto *dto.MenuDto, iconName string) []string {
	description := menuDto.Description
	if strings.TrimSpace(description) == "" {
		description = " "
	}
	tags := []string{
		fmt.Sprintf("showName:%s", menuDto.ShowName),
		"menu", // 标识为菜单
		fmt.Sprintf("description:%s", description),
		"parentName:menu.tag.appspace",
		fmt.Sprintf("openType:%d", menuDto.OpenType),
	}
	if iconName != "" {
		tags = append(tags, fmt.Sprintf("iconUrl:%s", iconName))
	}
	if len(menuDto.Tags) > 0 {
		tags = append(tags, menuDto.Tags...)
	}
	return tags
}

// fetchOrCreateService 根据 baseURL 获取或创建服务
func (l *MenuLogic) fetchOrCreateService(serviceName string, updateService bool, protocol, host string, port int) (string, string, error) {
	var serviceID string

	if serviceName != "" {
		service, err := l.kongLogic.QueryServiceJsonById(serviceName)
		if err != nil {
			return "", "", errors.NewBuzError(500, "menu.save.failed")
		}

		if service == nil {
			newService, err := l.kongLogic.CreateService(&KongServiceRequest{
				Name:     serviceName,
				Protocol: protocol,
				Host:     host,
				Port:     port,
				Enabled:  true,
			})
			if err != nil {
				return "", "", errors.NewBuzError(500, "menu.save.failed")
			}
			serviceID = newService.ID
		} else {
			if updateService {
				if !strings.EqualFold(protocol, cast.ToString(service["protocol"])) ||
					!strings.EqualFold(host, cast.ToString(service["host"])) ||
					port != cast.ToInt(service["port"]) {
					// service变更
					service["protocol"] = protocol
					service["host"] = host
					service["port"] = port
					l.kongLogic.UpdateService(cast.ToString(service["id"]), service)
				}
			}

			serviceID = cast.ToString(service["id"])
		}
	} else {
		serviceName = fmt.Sprintf("%s-%s-%d", protocol, host, port)
		service, err := l.kongLogic.QueryServiceJsonById(serviceName)
		if err != nil {
			return "", "", errors.NewBuzError(500, "menu.save.failed")
		}

		if service == nil {
			newService, err := l.kongLogic.CreateService(&KongServiceRequest{
				Name:     serviceName,
				Protocol: protocol,
				Host:     host,
				Port:     port,
				Enabled:  true,
			})
			if err != nil {
				return "", "", errors.NewBuzError(500, "menu.save.failed")
			}
			serviceID = newService.ID
		} else {
			serviceID = cast.ToString(service["id"])
		}
	}

	return serviceID, serviceName, nil
}
