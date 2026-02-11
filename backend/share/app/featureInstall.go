package app

import (
	"backend/share/app/adapter"
	"backend/share/app/model"
	"backend/share/app/util"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func InstallFeature(fc *model.NewFeatureModel) error {
	if err := fc.Validate(); err != nil {
		return err
	}

	// 加载镜像
	switch fc.GetImageSource() {
	case "local":
		{
			if err := adapter.LoadImageFromLocal(fc.ImagePath); err != nil {
				return err
			}
		}
	case "remote":
		{
			if err := adapter.PullImageFromRemote(fc.ImageUrl); err != nil {
				return err
			}
		}
	}
	if fc.ComposeYaml != "" {
		// 部署docker compose文件
		defaultNetwork := "tier0_edge_network"
		if err := adapter.DeployComposeYaml(fc.ComposeYaml, defaultNetwork); err != nil {
			return err
		}
	}
	// 注册菜单
	if err := adapter.SaveMenu(fc.Menu); err != nil {
		return err
	}
	// 将应用配置写到本地文件,文件名为fc.Name.json
	return saveFeatureConfig(fc)
}

// saveFeatureConfig 将功能配置保存到本地 JSON 文件
func saveFeatureConfig(fc *model.NewFeatureModel) error {
	if fc == nil {
		return fmt.Errorf("功能配置不能为空")
	}

	// 1. 确保配置目录存在
	if err := os.MkdirAll(util.AppInstalledDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 2. 构建文件路径
	// 使用功能名称作为文件名，确保文件名安全
	safeFileName := makeFileNameSafe(fc.Name)
	configPath := filepath.Join(util.AppInstalledDir, safeFileName+".json")

	// 设置安装时间为当前时间
	fc.InstallTime = time.Now().Format("2006-01-02 15:04:05")
	// 3. 转换为 JSON
	configJSON, err := json.MarshalIndent(fc, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 4. 写入文件
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	// 5. 记录日志
	log.Printf("[app]: 功能配置已保存: %s, 安装时间: %s", fc.Name, fc.InstallTime)
	return nil
}

func makeFileNameSafe(fileName string) string {
	// 定义非法字符
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	safeName := fileName

	// 替换非法字符为下划线
	for _, char := range invalidChars {
		safeName = strings.ReplaceAll(safeName, char, "_")
	}

	// 移除首尾空格
	safeName = strings.TrimSpace(safeName)

	// 如果文件名为空，使用默认名称
	if safeName == "" {
		safeName = "unnamed-feature"
	}

	// 限制文件名长度
	if len(safeName) > 100 {
		safeName = safeName[:100]
	}

	return safeName
}
