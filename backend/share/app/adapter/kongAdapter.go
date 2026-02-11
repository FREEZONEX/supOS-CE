package adapter

import (
	"backend/internal/adapters/kong/logic"
	"backend/internal/repo/relationDB"
	"backend/share/app/model"
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"github.com/kong/go-kong/kong"
)

func RegisterRouter(appHost string, port int, stripPath bool) error {
	// 1. 初始化 Kong 客户端
	client, err := kong.NewClient(kong.String("http://kong:8001"), nil)
	if err != nil {
		log.Fatalf("无法连接到 Kong: %v", err)
	}

	// 定义服务名和路由路径
	serviceName := appHost
	routePath := fmt.Sprintf("/apps/%s", appHost)

	ctx := context.Background()

	// 2. 注册 Service
	// 注意：Host 应该是你后端容器或服务的实际地址
	newService := &kong.Service{
		Name:     kong.String(serviceName),
		Protocol: kong.String("http"),
		Host:     kong.String(appHost),
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
		routeName := appHost
		newRoute := &kong.Route{
			Name:      kong.String(routeName),
			Paths:     kong.StringSlice(routePath),
			Service:   service,
			StripPath: kong.Bool(stripPath),
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

func SaveMenu(menu *model.MenuModel) error {
	log.Printf("[app]: 开始保存菜单: %s, indexUrl: %s", menu.Name, menu.IndexUrl)

	// 1. 解析indexUrl，分离host和port
	host, port, err := parseIndexUrl(menu.IndexUrl)
	if err != nil {
		log.Printf("[app]: 解析indexUrl失败: %v", err)
		return errors.Parameter.WithMsg(fmt.Sprintf("parse indexUrl failed: %v", err))
	}

	log.Printf("[app]: 解析结果 - Host: %s, Port: %d", host, port)
	// 2. 调用RegisterRouter，注册kong的service和router
	if err := RegisterRouter(host, port, menu.StripPath); err != nil {
		log.Printf("[app]: 注册Kong路由失败: %v", err)
		return errors.Parameter.WithMsg(fmt.Sprintf("register route failed: %v", err))
	}

	log.Printf("[app]: Kong路由注册成功")

	// 3. 将菜单数据保存到数据库表supos_resource
	if err := saveMenuToDatabase(menu); err != nil {
		log.Printf("[app]: 保存菜单到数据库失败: %v", err)
		// 注意：这里不返回错误，因为Kong路由已经注册成功
		// 数据库保存失败不影响路由功能
	}

	log.Printf("[app]: 菜单保存完成: %s", menu.Name)
	return nil
}

// parseIndexUrl 解析indexUrl，分离host和port
func parseIndexUrl(indexUrl string) (string, int, error) {
	if indexUrl == "" {
		return "", 0, fmt.Errorf("indexUrl不能为空")
	}

	// 移除协议前缀
	url := indexUrl
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	}

	// 分割host和port
	var host string
	var port int = 80 // 默认端口

	if strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		if len(parts) != 2 {
			return "", 0, fmt.Errorf("无效的URL格式: %s", indexUrl)
		}
		host = parts[0]
		portStr := parts[1]

		// 移除端口后的路径（如果有）
		if strings.Contains(portStr, "/") {
			portStr = strings.Split(portStr, "/")[0]
		}

		// 转换端口为整数
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			return "", 0, fmt.Errorf("无效的端口号: %s", portStr)
		}
		port = portInt
	} else {
		// 没有端口号，使用默认端口
		host = url
		// 移除路径部分（如果有）
		if strings.Contains(host, "/") {
			host = strings.Split(host, "/")[0]
		}
	}

	// 验证host
	if host == "" {
		return "", 0, fmt.Errorf("无法解析host: %s", indexUrl)
	}

	return host, port, nil
}

// saveMenuToDatabase 将菜单数据保存到数据库表supos_resource
func saveMenuToDatabase(menu *model.MenuModel) error {
	log.Printf("[app]: 开始保存菜单到数据库: %s", menu.Name)

	// 创建上下文
	ctx := context.Background()

	// 获取数据库连接
	db := relationDB.GetDb(ctx)
	if db == nil {
		log.Printf("[app]: 无法获取数据库连接")
		return errors.Parameter.WithMsg(fmt.Sprintf("can't connect to postgres database"))
	}

	// 检查是否已存在相同code的菜单
	var existingCount int64
	err := db.Model(&relationDB.SuposResource{}).
		Where("code = ?", menu.Name).
		Count(&existingCount).Error
	if err != nil {
		log.Printf("[app]: 查询现有菜单失败: %v", err)
		return errors.Parameter.WithMsg(fmt.Sprintf("search menu failed: %v", err))
	}

	u, err := url.Parse(menu.IndexUrl)
	if err != nil {
		panic(err)
	}

	serviceName := u.Hostname() // frontend-app
	path := u.Path

	// 父级菜单
	var p int64 = 4
	// 创建菜单资源对象
	menuResource := &relationDB.SuposResource{
		ParentID:        &p,
		Code:            serviceName,
		NameCode:        stringPtr(serviceName),
		URL:             stringPtr(path), // 相对路由
		DescriptionCode: stringPtr(menu.Description),
		URLType:         intPtr(2),
		Type:            2,              // 假设1表示菜单类型
		OpenType:        intPtr(0),      // 默认打开类型
		Sort:            intPtr(0),      // 默认排序
		Enable:          boolPtr(true),  // 启用
		EditEnable:      boolPtr(true),  // 可编辑
		HomeEnable:      boolPtr(true),  // 非首页
		Fixed:           boolPtr(false), // 非固定
		Icon:            stringPtr(menu.IconUrl),
		RouteSource:     intPtr(1),
		CreateAt:        time.Now(),
		UpdateAt:        time.Now(),
	}

	log.Printf("[app]: 菜单资源对象创建成功: Code=%s, URL=%s", menuResource.Code, *menuResource.URL)

	// 保存到数据库
	if existingCount > 0 {
		// 更新现有记录
		err = db.Model(&relationDB.SuposResource{}).
			Where("code = ?", menu.Name).
			Updates(menuResource).Error
		if err != nil {
			log.Printf("[app]: 更新菜单到数据库失败: %v", err)
			return errors.Parameter.WithMsg(fmt.Sprintf("update menu failed: %v", err))
		}
		log.Printf("[app]: 菜单更新到数据库成功: %s", menu.Name)
	} else {
		// 创建新记录
		err = db.Create(menuResource).Error
		if err != nil {
			log.Printf("[app]: 保存菜单到数据库失败: %v", err)
			return errors.Parameter.WithMsg(fmt.Sprintf("create menu failed: %v", err))
		}
		log.Printf("[app]: 菜单保存到数据库成功: %s (ID: %d)", menu.Name, menuResource.ID)
	}

	return nil
}

func deleteMenu(appName string) {
	// 创建 MenuLogic 实例
	menuLogic := logic.NewMenuLogic("kong", 8001)

	// 删除菜单
	err := menuLogic.DeleteMenu(appName)
	if err != nil {
		log.Printf("删除菜单失败: %v", err)
	} else {
		log.Println("菜单删除成功")
	}
}

// Helper functions for pointer creation
func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
