package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/errors"
	"backend/internal/common/event"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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
	s := &DashboardService{
		logger:       logx.WithContext(context.Background()),
		fileRootPath: constants.RootPath, // 暂定根路径，后续可从配置中读取
	}
	// Note: InitDashboardsOnStartup should be called by the application's main startup logic,
	// after the database connection is established and passed in.
	return s
}

// InitDashboardsOnStartup 应用启动时初始化 Dashboard
func (s *DashboardService) InitDashboardsOnStartup(ctx context.Context) {
	go func() {
		dashboardMapper := relationDB.DashboardMapper{}
		db := relationDB.GetDb(ctx)
		dashboards, err := dashboardMapper.SelectDashboardsToInit(db)
		if err != nil {
			s.logger.Errorf("failed to select dashboards to init: %v", err)
			return
		}

		if len(dashboards) == 0 {
			s.logger.Info("no dashboards to initialize")
			return
		}

		s.logger.Infof("dashboards to initialize: %d", len(dashboards))
		for _, dashboard := range dashboards {
			if dashboard.JsonContent == "" {
				continue
			}

			var dashboardData map[string]any
			if err := json.Unmarshal([]byte(dashboard.JsonContent), &dashboardData); err != nil {
				s.logger.Errorf("failed to unmarshal dashboard json content for %s: %v", dashboard.Name, err)
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
			existing, _ := grafanautil.GetDashboardByUUID(ctx, uid)
			if existing != nil {
				dashboard.NeedInit = false
				if err := dashboardMapper.UpdateById(db, dashboard); err != nil {
					s.logger.Errorf("failed to update dashboard init status for %s: %v", dashboard.Name, err)
				}
				s.logger.Infof("dashboard %s already initialized", dashboard.Name)
				continue
			}

			// 不存在，则创建
			dashboardMap["id"] = nil
			jsonBytes, _ := json.Marshal(dashboardData)
			_, err = grafanautil.CreateDashboardByBody(ctx, uid, "", string(jsonBytes))
			if err != nil {
				s.logger.Errorf("failed to initialize dashboard %s: %v", dashboard.Name, err)
			} else {
				dashboard.NeedInit = false
				if err := dashboardMapper.UpdateById(db, dashboard); err != nil {
					s.logger.Errorf("failed to update dashboard init status for %s after creation: %v", dashboard.Name, err)
				}
				s.logger.Infof("dashboard %s initialized successfully", dashboard.Name)
			}
		}
	}()
}

// OnEventRemoveTopics 当 UNS Topic 被删除时的处理逻辑
func (s *DashboardService) OnEventRemoveTopics(event *event.RemoveTopicsEvent) error {
	if event == nil {
		return errors.NewBuzError(400, "global.event.nil")
	}
	aliasList := base.Map(event.Topics, func(e *types.CreateTopicDto) string {
		return e.GetAlias()
	})
	return s.RemoveByUnsAliasList(aliasList)
}

func (s *DashboardService) RemoveByUnsAliasList(aliasList []string) error {
	s.logger.Infof("removing dashboards for topics: %v", aliasList)
	dashboardRefMapper := relationDB.DashboardRefMapper{}
	dashboardMapper := relationDB.DashboardMapper{}
	db := relationDB.GetDb(context.Background())
	// 1. 根据别名查询关联的 dashboard ID
	refs, err := dashboardRefMapper.SelectByUnsAliases(db, aliasList)
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
	_ = dashboardRefMapper.DeleteByUnsAliasList(db, aliasList)
	// 2. 批量删除 dashboard
	return dashboardMapper.DeleteBatchIds(db, idsToDelete)
}
func (s *DashboardService) CreateDashboards(ctx context.Context, dashboardVos []event.DashboardVo) error {
	now := time.Now()
	dashboards := make([]*relationDB.DashboardModel, 0, len(dashboardVos))
	refers := make([]*relationDB.DashboardRefModel, 0, len(dashboardVos))
	for _, dashboard := range dashboardVos {
		dashboards = append(dashboards, &relationDB.DashboardModel{
			ID:         dashboard.UUID,
			Name:       dashboard.Name,
			Creator:    dashboard.UserName,
			CreateTime: now,
			UpdateTime: now,
		})
		for _, unsAlias := range dashboard.UnsAlias {
			refers = append(refers, &relationDB.DashboardRefModel{
				DashboardID: dashboard.UUID,
				UnsAlias:    unsAlias,
			})
		}

	}
	db := relationDB.GetDb(ctx)
	dashboardMapper := relationDB.DashboardMapper{}
	err := dashboardMapper.SaveBatch(db, dashboards)
	if err != nil {
		s.logger.Errorf("failed to insert dashboard by event: %v", err)
		return err
	}
	// 创建引用关系
	dashboardRefMapper := relationDB.DashboardRefMapper{}
	return dashboardRefMapper.SaveBatch(db, refers)
}

// OnEventCreateDashboard 通过事件创建 Dashboard
func (s *DashboardService) OnEventCreateDashboard(event *event.CreateDashboardEvent) error {
	if event == nil {
		return errors.NewBuzError(400, "global.event.nil")
	}
	dashboardVos := event.Dashboards
	s.logger.Debugf("creating dashboard by event: %+v", dashboardVos)
	return s.CreateDashboards(event.Context, dashboardVos)
}

/*// DataExport 导出 Dashboard 数据
func (s *DashboardService) DataExport(ctx context.Context, db *gorm.DB, exportParam *dto.DashboardExportParam) (string, error) {
	exportCtx := &exporter.DashboardExportContext{}

	// 1. 获取数据
	err := s.fetchDataForExport(ctx, exportCtx, exportParam)
	if err != nil {
		s.logger.Errorf("failed to fetch data for export: %v", err)
		return "", errors.NewBuzError(500, "global.dashboard.export.error")
	}

	// 2. 导出数据到 JSON 文件
	exp := exporter.NewDashboardDataExporter()
	path, err := exp.ExportData(exportCtx, s.fileRootPath)
	if err != nil {
		s.logger.Errorf("failed to export data: %v", err)
		return "", errors.NewBuzError(500, "global.dashboard.export.error")
	}

	return path, nil
}

// fetchDataForExport 为导出获取数据
func (s *DashboardService) fetchDataForExport(ctx context.Context, dCtx *exporter.DashboardExportContext, exportParam *dto.DashboardExportParam) error {
	var dashboards []*relationDB.DashboardModel
	var err error
	db := relationDB.GetDb(ctx)
	var dashboardMapper relationDB.DashboardMapper
	if len(exportParam.Ids) > 0 {
		dashboards, err = dashboardMapper.SelectByIds(db, exportParam.Ids)
	} else if exportParam.ExportType == "ALL" {
		dashboards, err = dashboardMapper.SelectAll(db)
	}
	if err != nil {
		return err
	}

	for _, dashboard := range dashboards {
		if dashboard.Type == 1 { // Grafana
			jsonContent, err := grafanautil.Get(ctx, dashboard.ID)
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

	dCtx.DashboardModels = dashboards
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

	dashboardMapper := &relationDB.DashboardMapper{}
	importContext := importer.NewDashboardImportContext(filePath)
	dataImporter := importer.NewDashboardDataImporter(importContext, dashboardMapper)

	finalTask := "dashboard.create.task.name.final" // 假设从i18n获取

	if err := dataImporter.ImportData(ctx, file); err != nil {
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
*/
