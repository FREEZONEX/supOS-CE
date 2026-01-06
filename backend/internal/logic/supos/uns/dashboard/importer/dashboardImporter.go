package importer

import (
	"backend/internal/common/utils/fuxautil"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/repo/relationDB"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// DashboardImportContext Dashboard 导入上下文
type DashboardImportContext struct {
	FilePath      string
	Total         int
	CheckErrorMap map[string]string // key: dashboardID, value: error message
}

// NewDashboardImportContext 创建 DashboardImportContext 实例
func NewDashboardImportContext(filePath string) *DashboardImportContext {
	return &DashboardImportContext{
		FilePath:      filePath,
		CheckErrorMap: make(map[string]string),
	}
}

// AddError 添加错误信息
func (c *DashboardImportContext) AddError(id, errMsg string) {
	c.CheckErrorMap[id] = errMsg
}

// DataEmpty 检查是否有数据
func (c *DashboardImportContext) DataEmpty() bool {
	return c.Total == 0
}

// DashboardJsonWrapper Dashboard JSON 包装器
type DashboardJsonWrapper struct {
	Data []*relationDB.DashboardModel `json:"data"`
}

// DashboardDataImporter Dashboard 数据导入器
type DashboardDataImporter struct {
	context              *DashboardImportContext
	dashboardMapper      *relationDB.DashboardMapper
	dashboardJsonWrapper *DashboardJsonWrapper
	logger               logx.Logger
}

// NewDashboardDataImporter 创建 DashboardDataImporter 实例
func NewDashboardDataImporter(
	ctx *DashboardImportContext,
	dashboardMapper *relationDB.DashboardMapper,
) *DashboardDataImporter {
	return &DashboardDataImporter{
		context:         ctx,
		dashboardMapper: dashboardMapper,
		logger:          logx.WithContext(context.Background()),
	}
}

// ImportData 导入 Dashboard 数据
func (i *DashboardDataImporter) ImportData(ctx context.Context, file *os.File) error {
	// 解析 JSON 文件
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&i.dashboardJsonWrapper); err != nil {
		i.logger.Errorf("failed to parse json file: %v", err)
		return fmt.Errorf("dashboard.import.json.error")
	}

	// 处理导入数据
	return i.handleImportData(ctx)
}

// handleImportData 处理导入数据
func (i *DashboardDataImporter) handleImportData(ctx context.Context) error {
	if i.dashboardJsonWrapper == nil || len(i.dashboardJsonWrapper.Data) == 0 {
		return nil
	}

	i.context.Total = len(i.dashboardJsonWrapper.Data)
	now := time.Now()

	// 收集 ID 和名称
	ids := make([]string, 0, len(i.dashboardJsonWrapper.Data))
	names := make([]string, 0, len(i.dashboardJsonWrapper.Data))
	for _, dashboard := range i.dashboardJsonWrapper.Data {
		dashboard.CreateTime = now
		dashboard.UpdateTime = now
		ids = append(ids, dashboard.ID)
		names = append(names, dashboard.Name)
	}
	db := relationDB.GetDb(context.Background())
	// 检查 ID 是否已存在
	existByID := make(map[string]bool)
	for _, id := range ids {
		existing, err := i.dashboardMapper.SelectById(db, id)
		if err != nil {
			return err
		}
		if existing != nil {
			existByID[id] = true
		}
	}

	// 检查名称是否已存在
	existingByNames, err := i.dashboardMapper.SelectByFlowNames(db, names)
	if err != nil {
		return err
	}
	existByName := make(map[string]bool)
	for _, dashboard := range existingByNames {
		existByName[dashboard.Name] = true
	}

	// 筛选可以添加的数据
	addList := make([]*relationDB.DashboardModel, 0)
	for _, dashboard := range i.dashboardJsonWrapper.Data {
		if existByID[dashboard.ID] {
			i.context.AddError(dashboard.ID, "dashboard.id.already.exists")
		} else if existByName[dashboard.Name] {
			i.context.AddError(dashboard.ID, "dashboard.name.already.exists")
		} else {
			addList = append(addList, dashboard)
		}
	}

	// 批量插入数据
	if len(addList) > 0 {
		for _, dashboard := range addList {
			if err := i.dashboardMapper.Insert(db, dashboard); err != nil {
				i.logger.Errorf("failed to insert dashboard: %v", err)
				i.context.AddError(dashboard.ID, err.Error())
				continue
			}

			// 根据类型调用 Grafana/Fuxa API 创建 Dashboard
			if dashboard.Type == 1 {
				// Grafana Dashboard
				_, err := grafanautil.CreateDashboardByBody(ctx, dashboard.ID, "", dashboard.JsonContent)
				if err != nil {
					i.logger.Errorf("failed to create grafana dashboard: %v", err)
				}
			} else if dashboard.Type == 2 {
				// Fuxa Dashboard
				_, err := fuxautil.Create(dashboard.JsonContent)
				if err != nil {
					i.logger.Errorf("failed to create fuxa dashboard: %v", err)
				}
			}
		}
	}

	return nil
}

// WriteError 将错误信息写入文件
func (i *DashboardDataImporter) WriteError(fileRootPath string) (string, error) {
	// 构建错误文件路径
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("err_%s", filepath.Base(i.context.FilePath))
	relativePath := fmt.Sprintf("import/errors/dashboard/%s/%s", timestamp, filename)
	targetPath := filepath.Join(fileRootPath, relativePath)

	// 确保目录存在
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		i.logger.Errorf("failed to create error directory: %v", err)
		return "", err
	}

	// 为每个 Dashboard 添加错误信息
	for _, dashboard := range i.dashboardJsonWrapper.Data {
		if errMsg, ok := i.context.CheckErrorMap[dashboard.ID]; ok {
			dashboard.Error = errMsg
		}
	}

	// 构建错误数据结构
	errorData := map[string]interface{}{
		"data": i.dashboardJsonWrapper.Data,
	}

	// 创建文件
	file, err := os.Create(targetPath)
	if err != nil {
		i.logger.Errorf("failed to create error file: %v", err)
		return "", err
	}
	defer file.Close()

	// 写入 JSON 数据（格式化输出）
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(errorData); err != nil {
		i.logger.Errorf("failed to encode error data: %v", err)
		return "", err
	}

	i.logger.Infof("dashboard import error file created: %s", targetPath)
	return relativePath, nil
}
