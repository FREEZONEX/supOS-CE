package adapter

import (
	"backend/internal/adapters/kong/dto"
	"backend/internal/adapters/kong/logic"
	"backend/share/app/model"
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/kong/go-kong/kong"
)

func RegisterRouter(vendorName string, appName string, host string, port int) error {
	// 1. 初始化 Kong 客户端
	client, err := kong.NewClient(kong.String("http://kong:8001"), nil)
	if err != nil {
		log.Fatalf("无法连接到 Kong: %v", err)
	}

	// 定义服务名和路由路径
	serviceName := fmt.Sprintf("%s-%s-service", vendorName, appName)
	routePath := fmt.Sprintf("/apps/%s-%s", vendorName, appName)

	ctx := context.Background()

	// 2. 注册 Service
	// 注意：Host 应该是你后端容器或服务的实际地址
	newService := &kong.Service{
		Name:     kong.String(serviceName),
		Protocol: kong.String("http"),
		Host:     kong.String(host),
		Port:     kong.Int(port),
	}

	// 判断service是否存在
	_, err0 := client.Services.Get(ctx, kong.String(serviceName))
	if err0 != nil {
		// 如果服务不存在，创建新服务
		service, err := client.Services.Create(ctx, newService)
		if err != nil {
			log.Fatalf("创建 Service 失败: %v", err)
			return err
		}
		// 3. 注册 Route
		routeName := fmt.Sprintf("%s-%s-route", vendorName, appName)
		newRoute := &kong.Route{
			Name:  kong.String(routeName),
			Paths: kong.StringSlice(routePath),
			// 将此路由绑定到刚刚创建的 Service
			Service: service,
			// 是否在转发时剥离匹配的路径前缀
			StripPath: kong.Bool(true),
		}

		// 判断route是否存在
		_, err0 := client.Routes.Get(ctx, kong.String(routeName))
		if err0 != nil {
			// 如果路由不存在，创建新路由
			route, err := client.Routes.Create(ctx, newRoute)
			if err != nil {
				log.Fatalf("创建 Route 失败: %v", err)
			}
			log.Printf("Route 已注册: %s, 路径: %v", *route.Name, route.Paths)
		}
	}
	return nil

}

func SaveMenu(showName string, desc string, indexUrl string) error {
	// 从 indexUrl 中提取 IP 和端口
	ip, port, err := extractHostAndPort(indexUrl)
	if err != nil {
		return fmt.Errorf("解析 indexUrl 失败: %v", err)
	}

	// 构建服务名称（使用域名格式）
	appDomain := fmt.Sprintf("%s-%d", ip, port)
	// 替换点号为横杠，使其符合域名格式
	appDomain = strings.ReplaceAll(appDomain, ".", "-")

	menuDto := &dto.MenuDto{
		ServiceName: appDomain,
		Name:        appDomain,
		ShowName:    showName,
		Description: desc,
		BaseURL:     indexUrl,
		IsMenu:      true,
	}
	// 创建 MenuLogic 实例
	menuLogic := logic.NewMenuLogic("kong", 8001)
	// 方式1：完整创建菜单
	return menuLogic.CreateMenu(menuDto, true)
}

// extractHostAndPort 从 URL 中提取主机和端口
func extractHostAndPort(urlStr string) (string, int, error) {
	// 确保 URL 有协议前缀
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "http://" + urlStr
	}

	// 解析 URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", 0, fmt.Errorf("URL 解析失败: %v", err)
	}

	// 提取主机和端口
	host := parsedURL.Hostname()
	if host == "" {
		return "", 0, fmt.Errorf("URL 中未找到主机名")
	}

	portStr := parsedURL.Port()
	port := 80 // 默认 HTTP 端口

	if portStr != "" {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			return "", 0, fmt.Errorf("端口解析失败: %v", err)
		}
		port = portInt
	} else {
		// 根据协议设置默认端口
		if parsedURL.Scheme == "https" {
			port = 443
		}
	}

	return host, port, nil
}

func deleteMenu(appConf *model.AppConfig) {
	// 创建 MenuLogic 实例
	menuLogic := logic.NewMenuLogic("kong", 8001)

	// 删除菜单
	err := menuLogic.DeleteMenu(fmt.Sprintf("%s-%s", appConf.VendorName, appConf.Name))
	if err != nil {
		log.Printf("删除菜单失败: %v", err)
	} else {
		log.Println("菜单删除成功")
	}
}
