package app

import (
	"backend/share/app/model"
	"backend/share/app/util"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ListInstalledFeatures 查询已安装的功能列表
func ListInstalledFeatures() ([]model.NewFeatureModel, error) {
	// 1. 确保目录存在
	if err := os.MkdirAll(util.AppInstalledDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 2. 读取目录内容
	entries, err := os.ReadDir(util.AppInstalledDir)
	if err != nil {
		return nil, fmt.Errorf("读取安装目录失败: %v", err)
	}

	// 3. 初始化结果列表
	var features []model.NewFeatureModel

	// 4. 遍历目录中的文件
	for _, entry := range entries {
		// 跳过目录，只处理文件
		if entry.IsDir() {
			continue
		}

		// 只处理 .json 文件
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		// 5. 构建完整文件路径
		filePath := filepath.Join(util.AppInstalledDir, entry.Name())

		// 6. 读取并解析配置文件
		feature, err := loadFeatureFromFile(filePath)
		if err != nil {
			log.Printf("加载配置文件 %s 失败: %v", entry.Name(), err)
			continue
		}

		// 7. 添加到结果列表
		features = append(features, *feature)
	}

	// 8. 记录日志
	log.Printf("已加载 %d 个已安装功能", len(features))

	return features, nil
}

// loadFeatureFromFile 从文件加载功能配置
func loadFeatureFromFile(filePath string) (*model.NewFeatureModel, error) {
	// 1. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", filePath)
	}

	// 2. 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 3. 解析 JSON
	var feature model.NewFeatureModel
	if err := json.Unmarshal(data, &feature); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 4. 验证配置
	if err := feature.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	// 5. 记录日志
	log.Printf("成功加载功能: %s", feature.Name)

	return &feature, nil
}

// GetInstalledFeature 根据名称获取已安装的功能
func GetInstalledFeature(featureName string) (*model.NewFeatureModel, error) {
	// 1. 构建文件路径
	safeFileName := makeFileNameSafe(featureName)
	filePath := filepath.Join(util.AppInstalledDir, safeFileName+".json")

	// 2. 加载功能配置
	return loadFeatureFromFile(filePath)
}

// DeleteInstalledFeature 删除已安装的功能
func DeleteInstalledFeature(featureName string) error {
	// 1. 构建文件路径
	safeFileName := makeFileNameSafe(featureName)
	filePath := filepath.Join(util.AppInstalledDir, safeFileName+".json")

	// 2. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("功能 %s 未安装", featureName)
	}

	// 3. 删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除配置文件失败: %v", err)
	}

	// 4. 记录日志
	log.Printf("功能 %s 已删除", featureName)

	return nil
}

// UpdateInstalledFeature 更新已安装的功能
func UpdateInstalledFeature(feature *model.NewFeatureModel) error {
	// 1. 验证配置
	if err := feature.Validate(); err != nil {
		return err
	}

	// 2. 保存配置（会覆盖原有文件）
	return saveFeatureConfig(feature)
}

// SearchInstalledFeatures 搜索已安装的功能
func SearchInstalledFeatures(keyword string) ([]model.NewFeatureModel, error) {
	// 1. 获取所有功能
	allFeatures, err := ListInstalledFeatures()
	if err != nil {
		return nil, err
	}

	// 2. 如果没有关键词，返回所有功能
	if keyword == "" {
		return allFeatures, nil
	}

	// 3. 过滤功能
	var filteredFeatures []model.NewFeatureModel
	keyword = strings.ToLower(keyword)

	for _, feature := range allFeatures {
		// 在名称和描述中搜索
		if strings.Contains(strings.ToLower(feature.Name), keyword) ||
			strings.Contains(strings.ToLower(feature.Description), keyword) {
			filteredFeatures = append(filteredFeatures, feature)
		}
	}

	return filteredFeatures, nil
}
