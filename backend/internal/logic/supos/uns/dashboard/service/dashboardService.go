package service

import (
	"backend/internal/common/dto"
	"backend/internal/common/errors"
	"backend/internal/common/utils/fuxautil"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/logic/supos/uns/dashboard/dao"
	"backend/internal/logic/supos/uns/dashboard/exporter"
	"backend/internal/logic/supos/uns/dashboard/importer"
	"backend/internal/logic/supos/uns/dashboard/model"
	"backend/share/spring"
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// DashboardService Dashboard 业务逻辑 - 主要负责启动初始化、事件处理和导入导出等后台任务
type DashboardService struct {
	ctx          context.Context
	logger       logx.Logger
	fileRootPath string // 文件根路径，用于导入导出
}

func init() {
	spring.RegisterBean(NewDashboardService())
}

// NewDashboardService 创建 DashboardService 实例
func NewDashboardService() *DashboardService {
	ctx := context.Background()

	s := &DashboardService{
		ctx:          ctx,
		logger:       logx.WithContext(ctx),
		fileRootPath: "data", // 暂定根路径，后续可从配置中读取
	}
	// Note: InitDashboardsOnStartup should be called by the application's main startup logic,
	// after the database connection is established and passed in.
	return s
}

// InitDashboardsOnStartup 应用启动时初始化 Dashboard
func (s *DashboardService) InitDashboardsOnStartup(db *gorm.DB) {
	go func() {
		dashboardMapper := dao.NewDashboardMapper(db, s.ctx)
		dashboards, err := dashboardMapper.SelectDashboardsToInit()
		if err != nil {
			s.logger.Errorf("failed to select dashboards to init: %v", err)
			return
		}

		if len(dashboards) == 0 {
			s.logger.Info("no dashboards to initialize")
			return
		}

		s.logger.Infof("dashboards to initialize: %d", len(dashboards))
		for _, db := range dashboards {
			if db.JsonContent == "" {
				continue
			}

			var dashboardData map[string]any
			if err := json.Unmarshal([]byte(db.JsonContent), &dashboardData); err != nil {
				s.logger.Errorf("failed to unmarshal dashboard json content for %s: %v", db.Name, err)
				continue
			}

			dashboardMap, ok := dashboardData["dashboard"].(map[string]any)
			if !ok {
				continue
			}

			uid, ok := dashboardMap["uid"].(string)
			if !ok || uid == "" {
				continue
			}

			// 检查 Grafana 中是否已存在
			existing, _ := grafanautil.GetDashboardByUUID(uid)
			if existing != nil {
				db.NeedInit = false
				if err := dashboardMapper.UpdateById(db); err != nil {
					s.logger.Errorf("failed to update dashboard init status for %s: %v", db.Name, err)
				}
				s.logger.Infof("dashboard %s already initialized", db.Name)
				continue
			}

			// 不存在，则创建
			dashboardMap["id"] = nil
			jsonBytes, _ := json.Marshal(dashboardData)
			_, err = grafanautil.CreateDashboardByBody(uid, "", string(jsonBytes))
			if err != nil {
				s.logger.Errorf("failed to initialize dashboard %s: %v", db.Name, err)
			} else {
				db.NeedInit = false
				if err := dashboardMapper.UpdateById(db); err != nil {
					s.logger.Errorf("failed to update dashboard init status for %s after creation: %v", db.Name, err)
				}
				s.logger.Infof("dashboard %s initialized successfully", db.Name)
			}
		}
	}()
}

// OnRemoveTopics 当 UNS Topic 被删除时的处理逻辑
func (s *DashboardService) OnRemoveTopics(ctx context.Context, db *gorm.DB, aliases []string) error {
	s.logger.Infof("removing dashboards for topics: %v", aliases)
	dashboardRefMapper := dao.NewDashboardRefMapper(db, ctx)
	dashboardMapper := dao.NewDashboardMapper(db, ctx)

	// 1. 根据别名查询关联的 dashboard ID
	refs, err := dashboardRefMapper.SelectByUnsAliases(aliases)
	if err != nil {
		s.logger.Errorf("failed to select dashboard refs by aliases: %v", err)
		return err
	}
	if len(refs) == 0 {
		return nil
	}

	idsToDelete := make([]string, len(refs))
	for i, ref := range refs {
		idsToDelete[i] = ref.DashboardID
	}

	// 2. 批量删除 dashboard
	return dashboardMapper.DeleteBatchIds(idsToDelete)
}

// CreateByEvent 通过事件创建 Dashboard
func (s *DashboardService) CreateByEvent(ctx context.Context, db *gorm.DB, uuid, name, username string) error {
	s.logger.Infof("creating dashboard by event: name=%s, uuid=%s", name, uuid)
	now := time.Now()
	dashboard := &model.DashboardModel{
		ID:         uuid,
		Name:       name,
		Creator:    username,
		CreateTime: now,
		UpdateTime: now,
	}

	dashboardMapper := dao.NewDashboardMapper(db, ctx)
	err := dashboardMapper.Insert(dashboard)
	if err != nil {
		s.logger.Errorf("failed to insert dashboard by event: %v", err)
		return err
	}

	// 创建引用关系
	ref := &model.DashboardRefModel{
		DashboardID: uuid,
		UnsAlias:    name, // 假设 alias 和 name 相同
		CreateAt:    now,
	}
	dashboardRefMapper := dao.NewDashboardRefMapper(db, ctx)
	return dashboardRefMapper.Insert(ref)
}

// DataExport 导出 Dashboard 数据
func (s *DashboardService) DataExport(ctx context.Context, db *gorm.DB, exportParam *dto.DashboardExportParam) (string, error) {
	context := &exporter.DashboardExportContext{}
	dashboardMapper := dao.NewDashboardMapper(db, ctx)

	// 1. 获取数据
	err := s.fetchDataForExport(context, dashboardMapper, exportParam)
	if err != nil {
		s.logger.Errorf("failed to fetch data for export: %v", err)
		return "", errors.NewBuzError(500, "global.dashboard.export.error")
	}

	// 2. 导出数据到 JSON 文件
	exp := exporter.NewDashboardDataExporter()
	path, err := exp.ExportData(context, s.fileRootPath)
	if err != nil {
		s.logger.Errorf("failed to export data: %v", err)
		return "", errors.NewBuzError(500, "global.dashboard.export.error")
	}

	return path, nil
}

// fetchDataForExport 为导出获取数据
func (s *DashboardService) fetchDataForExport(context *exporter.DashboardExportContext, dashboardMapper *dao.DashboardMapper, exportParam *dto.DashboardExportParam) error {
	var dashboards []*model.DashboardModel
	var err error

	if len(exportParam.Ids) > 0 {
		dashboards, err = dashboardMapper.SelectByIds(exportParam.Ids)
	} else if exportParam.ExportType == "ALL" {
		dashboards, err = dashboardMapper.SelectAll()
	}
	if err != nil {
		return err
	}

	for _, dashboard := range dashboards {
		if dashboard.Type == 1 { // Grafana
			jsonContent, err := grafanautil.Get(dashboard.ID)
			if err != nil {
				s.logger.Errorf("failed to get grafana dashboard content for %s: %v", dashboard.ID, err)
			} else {
				dashboard.JsonContent = jsonContent
			}
		} else if dashboard.Type == 2 { // Fuxa
			jsonContent, err := fuxautil.Get(dashboard.ID)
			if err != nil {
				s.logger.Errorf("failed to get fuxa dashboard content for %s: %v", dashboard.ID, err)
			} else {
				dashboard.JsonContent = jsonContent
			}
		}
		dashboard.CreateTime = time.Time{} // 导出时清空时间
		dashboard.UpdateTime = time.Time{}
	}

	context.DashboardModels = dashboards
	return nil
}

// AsyncImport 异步导入 Dashboard 数据
func (s *DashboardService) AsyncImport(ctx context.Context, db *gorm.DB, filePath string) (*dto.RunningStatus, error) {
	file, err := os.Open(filePath)
	if err != nil {
		s.logger.Errorf("failed to open import file %s: %v", filePath, err)
		return dto.NewRunningStatus(400, "global.import.file.not.exist"), nil
	}
	defer file.Close()

	dashboardMapper := dao.NewDashboardMapper(db, ctx)
	importContext := importer.NewDashboardImportContext(filePath)
	dataImporter := importer.NewDashboardDataImporter(importContext, dashboardMapper)

	finalTask := "dashboard.create.task.name.final" // 假设从i18n获取

	if err := dataImporter.ImportData(file); err != nil {
		s.logger.Errorf("failed to import data from %s: %v", filePath, err)
		// 导入失败，尝试写入错误文件
		_, writeErr := s.writeImportErrorFile(dataImporter)
		if writeErr != nil {
			s.logger.Errorf("failed to write error file after import failure: %v", writeErr)
		}
		return dto.NewRunningStatus(500, err.Error()).SetTask(finalTask).SetProgress(0.0), nil
	}

	if importContext.DataEmpty() {
		return dto.NewRunningStatus(400, "dashboard.import.excel.empty"), nil
	}

	if len(importContext.CheckErrorMap) == 0 {
		status := dto.NewRunningStatus(200, "dashboard.import.rs.ok").SetTask(finalTask).SetProgress(100.0)
		status.TotalCount = importContext.Total
		status.SuccessCount = importContext.Total
		status.ErrorCount = 0
		return status, nil
	}

	// 存在部分错误
	relativePath, err := s.writeImportErrorFile(dataImporter)
	if err != nil {
		s.logger.Errorf("failed to write import error file: %v", err)
		// 即使写入错误文件失败，也要通知前端导入已完成但有错误
		status := dto.NewRunningStatus(206, "dashboard.import.rs.hasErr").SetTask(finalTask).SetProgress(100.0)
		status.TotalCount = importContext.Total
		status.ErrorCount = len(importContext.CheckErrorMap)
		status.SuccessCount = status.TotalCount - status.ErrorCount
		return status, nil
	}

	message := "dashboard.import.rs.hasErr"
	status := dto.NewRunningStatus(206, message, relativePath).SetTask(finalTask).SetProgress(100.0)
	status.TotalCount = importContext.Total
	status.ErrorCount = len(importContext.CheckErrorMap)
	status.SuccessCount = status.TotalCount - status.ErrorCount
	if status.SuccessCount == 0 {
		status.Msg = "global.import.rs.allErr"
	}
	return status, nil
}

func (s *DashboardService) writeImportErrorFile(dataImporter *importer.DashboardDataImporter) (string, error) {
	// 注意：这里的 outFile 参数在原始实现中是 *os.File，但在 Go 版本中，
	// 我们直接在 WriteError 内部处理文件创建，简化调用。
	// 这里我们假设 WriteError 内部处理文件创建。
	return dataImporter.WriteError(s.fileRootPath)
}
