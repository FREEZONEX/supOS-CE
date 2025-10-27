package listener

import (
	"os"
	"path/filepath"

	"github.com/magiconair/properties"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// LocalMenuCheckedStoragePath 本地菜单配置文件路径
	LocalMenuCheckedStoragePath = "/data/menu/check-menu.properties"
)

var (
	// localMenus 本地菜单存储，使用 concurrent-map 保证并发安全
	localMenus = cmap.New[string]()
)

// LoadLocalMenus 在服务启动时加载本地菜单配置
func LoadLocalMenus() {
	// 检查文件是否存在
	if _, err := os.Stat(LocalMenuCheckedStoragePath); os.IsNotExist(err) {
		logx.Info("==> skip load local menu data, because file is not exist")
		// 尝试创建目录
		if err := os.MkdirAll(filepath.Dir(LocalMenuCheckedStoragePath), 0755); err != nil {
			logx.Errorf("failed to create directory for local menu: %v", err)
		}
		return
	}

	// 使用 properties 库加载配置文件
	props, err := properties.LoadFile(LocalMenuCheckedStoragePath, properties.UTF8)
	if err != nil {
		logx.Errorf("local file(%s) load error: %v", LocalMenuCheckedStoragePath, err)
		return
	}

	// 将所有配置项存入 sync.Map
	for _, key := range props.Keys() {
		if value, ok := props.Get(key); ok {
			localMenus.Set(key, value)
		}
	}

	logx.Info("==> load success, checked menu loaded")
}

// GetLocalMenus 返回加载的本地菜单
func GetLocalMenus() map[string]string {
	return localMenus.Items()
}

// StoreLocalMenu 存储单个菜单项
func StoreLocalMenu(key, value string) {
	localMenus.Set(key, value)
}

// UpdateLocalMenus 批量更新本地菜单缓存
func UpdateLocalMenus(menus map[string]string) {
	// 清空原有数据
	localMenus.Clear()
	// 设置新数据
	for key, value := range menus {
		localMenus.Set(key, value)
	}
}
