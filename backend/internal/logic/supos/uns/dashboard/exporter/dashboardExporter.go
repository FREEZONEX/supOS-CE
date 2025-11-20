package exporter

import (
	"backend/internal/common/constants"
	"backend/internal/repo/relationDB"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// DashboardExportContext Dashboard 导出上下文
type DashboardExportContext struct {
	DashboardModels []*relationDB.DashboardModel
}

// DashboardDataExporter Dashboard 数据导出器
type DashboardDataExporter struct {
	logger logx.Logger
}

// NewDashboardDataExporter 创建 DashboardDataExporter 实例
func NewDashboardDataExporter() *DashboardDataExporter {
	return &DashboardDataExporter{
		logger: logx.WithContext(context.Background()),
	}
}

// ExportData 导出 Dashboard 数据到 JSON 文件
func (e *DashboardDataExporter) ExportData(context *DashboardExportContext, fileRootPath string) (string, error) {
	// 构建导出路径
	timestamp := time.Now().Format("20060102150405")
	relativePath := constants.GlobalExport + fmt.Sprintf("dashboard/%s/", timestamp) + constants.GlobalExportDashboard
	targetPath := filepath.Join(fileRootPath, relativePath)

	// 确保目录存在
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		e.logger.Errorf("failed to create export directory: %v", err)
		return "", err
	}

	// 构建导出数据结构
	exportData := map[string]interface{}{
		"data": context.DashboardModels,
	}

	// 创建文件
	file, err := os.Create(targetPath)
	if err != nil {
		e.logger.Errorf("failed to create export file: %v", err)
		return "", err
	}
	defer file.Close()

	// 写入 JSON 数据（格式化输出）
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportData); err != nil {
		e.logger.Errorf("failed to encode export data: %v", err)
		return "", err
	}

	e.logger.Infof("dashboard export success: %s", targetPath)
	return relativePath, nil
}
