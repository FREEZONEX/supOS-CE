package service

import (
	"backend/internal/common/dto"
	"backend/internal/common/errors"
	"backend/internal/common/utils/dbutil"
	"backend/internal/common/utils/fuxautil"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/logic/supos/uns/dashboard/dao"
	"backend/internal/logic/supos/uns/dashboard/exporter"
	"backend/internal/logic/supos/uns/dashboard/importer"
	"backend/internal/logic/supos/uns/dashboard/model"
	unsservice "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"os"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

// DashboardService Dashboard 业务逻辑
type DashboardService struct {
	ctx                 context.Context
	logger              logx.Logger
	dashboardMapper     *dao.DashboardMapper
	dashboardRefMapper  *dao.DashboardRefMapper
	dashboardMarkMapper *dao.DashboardMarkedMapper
	unsQueryService     *unsservice.UnsQueryService
	unsUpdateService    *unsservice.UnsUpdateService
	fileRootPath        string // 文件根路径，用于导入导出
}

// NewDashboardService 创建 DashboardLogic 实例
func NewDashboardService(
	ctx context.Context,
	dashboardMapper *dao.DashboardMapper,
	dashboardRefMapper *dao.DashboardRefMapper,
	dashboardMarkMapper *dao.DashboardMarkedMapper,
	unsQueryService *unsservice.UnsQueryService,
	unsUpdateService *unsservice.UnsUpdateService,
) *DashboardService {
	s := &DashboardService{
		ctx:                 ctx,
		logger:              logx.WithContext(ctx),
		dashboardMapper:     dashboardMapper,
		dashboardRefMapper:  dashboardRefMapper,
		dashboardMarkMapper: dashboardMarkMapper,
		unsQueryService:     unsQueryService,
		unsUpdateService:    unsUpdateService,
		fileRootPath:        "data", // 暂定根路径，后续可从配置中读取
	}
	s.InitDashboardsOnStartup() // 应用启动时初始化
	return s
}

// InitDashboardsOnStartup 应用启动时初始化 Dashboard
func (s *DashboardService) InitDashboardsOnStartup() {
	go func() {
		dashboards, err := s.dashboardMapper.SelectDashboardsToInit()
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
				if err := s.dashboardMapper.UpdateById(db); err != nil {
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
				if err := s.dashboardMapper.UpdateById(db); err != nil {
					s.logger.Errorf("failed to update dashboard init status for %s after creation: %v", db.Name, err)
				}
				s.logger.Infof("dashboard %s initialized successfully", db.Name)
			}
		}
	}()
}

// OnRemoveTopics 当 UNS Topic 被删除时的处理逻辑
func (s *DashboardService) OnRemoveTopics(aliases []string) error {
	s.logger.Infof("removing dashboards for topics: %v", aliases)
	// 1. 根据别名查询关联的 dashboard ID
	refs, err := s.dashboardRefMapper.SelectByUnsAliases(aliases)
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
	return s.dashboardMapper.DeleteBatchIds(idsToDelete)
}

// CreateByEvent 通过事件创建 Dashboard
func (s *DashboardService) CreateByEvent(uuid, name, username string) error {
	s.logger.Infof("creating dashboard by event: name=%s, uuid=%s", name, uuid)
	now := time.Now()
	dashboard := &model.DashboardModel{
		ID:         uuid,
		Name:       name,
		Creator:    username,
		CreateTime: now,
		UpdateTime: now,
	}

	err := s.dashboardMapper.Insert(dashboard)
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
	return s.dashboardRefMapper.Insert(ref)
}

// PageList 分页查询 Dashboard
func (s *DashboardService) PageList(
	keyword string,
	typ *int,
	orderCode string,
	descOrAsc string,
	pageNo int64,
	pageSize int64,
	userID string,
) (*dto.PageResultDTO[*model.DashboardModel], error) {
	// SQL 注入防护：转义关键字
	keyword = dbutil.EscapeForLike(keyword)

	// 排序字段校验
	if orderCode != "" {
		if orderCode != "name" && orderCode != "createTime" {
			return nil, errors.NewBuzError(400, "illegal sort param")
		}
		// 驼峰转下划线
		orderCode = camelToSnake(orderCode)
	}

	// 排序方向
	if descOrAsc == "" || descOrAsc != "ASC" {
		descOrAsc = "DESC"
	}

	// 查询总数
	total, err := s.dashboardMapper.SelectDashboardCount(keyword, typ)
	if err != nil {
		return nil, err
	}

	// 查询数据
	dashboards, err := s.dashboardMapper.SelectDashboard(userID, keyword, typ, orderCode, descOrAsc, pageNo, pageSize)
	if err != nil {
		return nil, err
	}

	// 构建响应
	data := make([]*model.DashboardModel, len(dashboards))
	for i, db := range dashboards {
		data[i] = &db.DashboardModel
	}

	return &dto.PageResultDTO[*model.DashboardModel]{
		Code:     0,
		Total:    total,
		PageNo:   pageNo,
		PageSize: pageSize,
		Data:     data,
	}, nil
}

// GetById 根据 ID 获取 Dashboard
func (s *DashboardService) GetById(id string) (*model.DashboardModel, error) {
	return s.dashboardMapper.SelectById(id)
}

// Create 创建 Dashboard
func (s *DashboardService) Create(dashboard *model.DashboardModel, creator string) (*model.DashboardModel, error) {
	// 检查名称是否重复
	dashboards, err := s.dashboardMapper.SelectByFlowNames([]string{dashboard.Name})
	if err != nil {
		return nil, err
	}
	if len(dashboards) > 0 {
		for _, db := range dashboards {
			if db.Type == dashboard.Type {
				return nil, errors.NewBuzError(500, "uns.dashboard.name.duplicate")
			}
		}
	}

	// 生成 ID
	dashboard.ID = uuid.New().String()
	dashboard.Creator = creator
	dashboard.CreateTime = time.Now()
	dashboard.UpdateTime = time.Now()

	// Grafana Dashboard 创建
	if dashboard.Type == 1 {
		// 构建 Dashboard JSON
		dashboardJSON := fmt.Sprintf(`{
			"dashboard": {
				"uid": "%s",
				"title": "%s",
				"id": null
			}
		}`, dashboard.ID, dashboard.Name)

		// 调用 Grafana API 创建 Dashboard
		url := grafanautil.GetGrafanaURL() + "/api/dashboards/db"
		_, err := grafanautil.CreateDashboardByBody(dashboard.ID, "", dashboardJSON)
		if err != nil {
			s.logger.Errorf("failed to create grafana dashboard: %v", err)
			return nil, errors.NewBuzError(500, "uns.dashboard.create.failed")
		}
		s.logger.Infof("created grafana dashboard: %s, url: %s", dashboard.ID, url)
	}

	// 保存到数据库
	err = s.dashboardMapper.Insert(dashboard)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

// Edit 编辑 Dashboard
func (s *DashboardService) Edit(dashboard *model.DashboardModel) error {
	// 检查 Dashboard 是否存在
	existing, err := s.dashboardMapper.SelectById(dashboard.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.NewBuzError(400, "uns.dashboard.not.exit")
	}

	// Grafana Dashboard 更新
	if existing.Type == 1 {
		// 获取现有的 Dashboard
		dbJSON, err := grafanautil.GetDashboardByUUID(dashboard.ID)
		if err != nil || dbJSON == nil {
			return errors.NewBuzError(400, "uns.dashboard.not.exit")
		}

		// 更新 title 和 description
		if dashboardObj, ok := dbJSON["dashboard"].(map[string]any); ok {
			dashboardObj["title"] = dashboard.Name
			dashboardObj["description"] = dashboard.Description
		}

		// 调用 Grafana API 更新
		jsonBytes, _ := json.Marshal(dbJSON)
		url := grafanautil.GetGrafanaURL() + "/api/dashboards/db"
		_, err = grafanautil.CreateDashboardByBody(dashboard.ID, "", string(jsonBytes))
		if err != nil {
			s.logger.Errorf("failed to update grafana dashboard: %v", err)
			return err
		}
		s.logger.Infof("updated grafana dashboard: %s, url: %s", dashboard.ID, url)
	}

	// 更新数据库
	dashboard.UpdateTime = time.Now()
	return s.dashboardMapper.UpdateById(dashboard)
}

// Delete 删除 Dashboard
func (s *DashboardService) Delete(uid string) error {
	// 检查 Dashboard 是否存在
	dashboard, err := s.dashboardMapper.SelectById(uid)
	if err != nil {
		return err
	}
	if dashboard == nil {
		return nil // 已经不存在，视为成功
	}

	// Grafana Dashboard 删除
	if dashboard.Type == 1 {
		err := grafanautil.DeleteDashboard(uid)
		if err != nil {
			s.logger.Errorf("failed to delete grafana dashboard: %v", err)
		}
	}

	// Fuxa Dashboard 删除
	if dashboard.Type == 2 {
		// Fuxa 使用 HTTP DELETE 请求删除
		url := fmt.Sprintf("%s/api/project/%s", fuxautil.GetFuxaURL(), uid)
		s.logger.Infof("deleting fuxa dashboard: %s", url)
		// 注意：fuxautil 目前没有 Delete 方法，需要直接 HTTP 调用或添加方法
	}

	// 删除置顶标记
	err = s.dashboardMarkMapper.DeleteById(uid)
	if err != nil {
		s.logger.Errorf("failed to delete dashboard mark: %v", err)
	}

	// 删除引用关系
	err = s.dashboardRefMapper.DeleteByDashboardId(uid)
	if err != nil {
		s.logger.Errorf("failed to delete dashboard ref: %v", err)
	}

	// 删除 Dashboard
	return s.dashboardMapper.DeleteById(uid)
}

// GetByUuid 根据 UUID 获取 Grafana Dashboard
func (s *DashboardService) GetByUuid(uuid string) (map[string]any, error) {
	dbJSON, err := grafanautil.GetDashboardByUUID(uuid)
	if err != nil || dbJSON == nil {
		return nil, errors.NewBuzError(400, "uns.dashboard.not.exit")
	}
	return dbJSON, nil
}

// IsExist 检查 Dashboard 是否存在
func (s *DashboardService) IsExist(alias string) (map[string]any, error) {
	uuid := grafanautil.GetDashboardUUIDByAlias(alias)
	dbJSON, err := grafanautil.GetDashboardByUUID(uuid)
	if err != nil || dbJSON == nil {
		return nil, errors.NewBuzError(400, "uns.dashboard.not.exit")
	}
	return dbJSON, nil
}

// MarkTop 置顶 Dashboard
func (s *DashboardService) MarkTop(id string, userID string) error {
	mark := &model.DashboardMarkModel{
		ID:     id,
		UserID: userID,
	}
	return s.dashboardMarkMapper.Insert(mark)
}

// RemoveMarkedTop 取消置顶 Dashboard
func (s *DashboardService) RemoveMarkedTop(id string, userID string) error {
	return s.dashboardMarkMapper.Delete(id, userID)
}

// BindUns 绑定 Dashboard 到 UNS
func (s *DashboardService) BindUns(dashboardID string, unsAlias string) error {
	// 检查 Dashboard 是否存在
	dashboard, err := s.dashboardMapper.SelectById(dashboardID)
	if err != nil {
		return err
	}
	if dashboard == nil {
		return errors.NewBuzError(400, "uns.dashboard.not.exit")
	}

	// 检查 UNS 是否存在
	unsResp, err := s.unsQueryService.GetModelDefinition(s.ctx, &types.ModelDetailReq{}, unsAlias)
	if err != nil {
		s.logger.Errorf("failed to get uns definition for alias %s: %v", unsAlias, err)
		return errors.NewBuzError(500, "uns.file.not.exist")
	}
	if unsResp == nil || unsResp.Data == nil || unsResp.Data.Id == "" {
		return errors.NewBuzError(400, "uns.file.not.exist")
	}

	// 删除旧的绑定关系
	err = s.dashboardRefMapper.DeleteByDashboardId(dashboardID)
	if err != nil {
		s.logger.Errorf("failed to delete old dashboard ref: %v", err)
	}

	// 创建新的绑定关系
	ref := &model.DashboardRefModel{
		DashboardID: dashboardID,
		UnsAlias:    unsAlias,
		CreateAt:    time.Now(),
	}
	return s.dashboardRefMapper.Insert(ref)
}

// GetByUns 根据 UNS 别名获取 Dashboard
func (s *DashboardService) GetByUns(unsAlias string) (*model.DashboardModel, error) {
	// TODO: 当前的 unsQueryService.GetModelDefinition 返回的 DTO 中不包含 Refers 字段，
	// 无法直接判断 UNS 是否为引用类型。此处的逻辑暂时简化为直接查询，
	// 后续如果 uns service 提供了更详细的接口，需要回来完善引用类型的处理逻辑。

	// unsResp, err := s.unsQueryService.GetModelDefinition(s.ctx, &types.ModelDetailReq{}, unsAlias)
	// if err != nil {
	// 	s.logger.Errorf("could not find uns definition for alias %s: %v", unsAlias, err)
	// 	// 即使找不到定义，也继续尝试查询，保持旧逻辑兼容性
	// 	return s.dashboardRefMapper.GetByUns(unsAlias)
	// }

	// if unsResp != nil && unsResp.Data != nil {
	// 	unsDef := unsResp.Data
	// 	// 检查是否是引用类型
	// 	if unsDef.DataType == constants.CitingType {
	// 		// DTO 中没有 Refers 字段，无法实现
	// 	}
	// }

	return s.dashboardRefMapper.GetByUns(unsAlias)
}

// CreateGrafanaByUns 基于 UNS 创建 Grafana Dashboard
func (s *DashboardService) CreateGrafanaByUns(alias string) (string, error) {
	// 1. 获取 UNS 定义
	unsResp, err := s.unsQueryService.GetModelDefinition(s.ctx, &types.ModelDetailReq{}, alias)
	if err != nil || unsResp == nil || unsResp.Data == nil {
		return "", errors.NewBuzError(400, "uns.file.not.exist")
	}
	unsDef := unsResp.Data

	// 2. 检查是否已有关联的 Dashboard
	existingDashboard, err := s.GetByUns(alias)
	if err != nil {
		s.logger.Errorf("error checking for existing dashboard for uns %s: %v", alias, err)
		// Fall through, but log the error
	}
	if existingDashboard != nil {
		return "", errors.NewBuzError(400, "uns.dashboard.already.exists")
	}

	// 3. 根据 UNS 字段构建 Grafana Dashboard JSON
	dashboardUID := uuid.New().String()
	dashboardJSON, err := s.buildGrafanaJSONFromUns(unsDef, dashboardUID)
	if err != nil {
		s.logger.Errorf("failed to build grafana json for uns %s: %v", alias, err)
		return "", errors.NewBuzError(500, "uns.dashboard.create.failed")
	}

	// 4. 调用 Grafana API 创建 Dashboard
	_, err = grafanautil.CreateDashboardByBody(dashboardUID, "", dashboardJSON)
	if err != nil {
		s.logger.Errorf("failed to create grafana dashboard for uns %s: %v", alias, err)
		return "", errors.NewBuzError(500, "uns.dashboard.create.failed")
	}

	// 5. 保存 Dashboard 记录到数据库
	now := time.Now()
	dashboard := &model.DashboardModel{
		ID:          dashboardUID,
		Name:        unsDef.Name,
		Creator:     "system", // 系统自动创建
		CreateTime:  now,
		UpdateTime:  now,
		Type:        1, // 1 for Grafana
		JsonContent: dashboardJSON,
		NeedInit:    false, // 已在 Grafana 中创建
	}
	if err = s.dashboardMapper.Insert(dashboard); err != nil {
		s.logger.Errorf("failed to save dashboard record for uns %s: %v", alias, err)
		// 尝试回滚 Grafana 的创建操作
		_ = grafanautil.DeleteDashboard(dashboardUID)
		return "", err
	}

	// 6. 创建 Dashboard 和 UNS 的绑定关系
	ref := &model.DashboardRefModel{
		DashboardID: dashboardUID,
		UnsAlias:    alias,
		CreateAt:    now,
	}
	if err = s.dashboardRefMapper.Insert(ref); err != nil {
		// 此处为非关键路径，只记录日志
		s.logger.Errorf("failed to bind dashboard %s to uns %s: %v", dashboardUID, alias, err)
	}

	// 7. 更新 UNS 的 flags
	addDashboard := true
	updateDto := &types.UpdateUnsDto{
		Alias:        alias,
		AddDashBoard: &addDashboard,
	}
	_, updateErr := s.unsUpdateService.UpdateDetail(s.ctx, updateDto)
	if updateErr != nil {
		// 此处为非关键路径，只记录日志
		s.logger.Errorf("failed to update uns flags for %s after creating dashboard: %v", alias, updateErr)
	}

	// 8. 返回新创建的 Dashboard UID
	return dashboardUID, nil
}

// buildGrafanaJSONFromUns 根据 UNS 定义构建 Grafana Dashboard JSON
func (s *DashboardService) buildGrafanaJSONFromUns(unsDef *types.ModelDetail, uid string) (string, error) {
	// 筛选出数值类型的字段用于创建图表
	numericFields := make([]*types.FieldDefine, 0)
	for _, field := range unsDef.Fields {
		if types.FieldType(field.Type).IsNumber() && !field.IsSystemField() {
			numericFields = append(numericFields, field)
		}
	}

	panels := make([]map[string]any, 0, len(numericFields))
	for i, field := range numericFields {
		panel := map[string]any{
			"id":    i + 1,
			"title": field.Name,
			"type":  "timeseries",
			"gridPos": map[string]int{
				"h": 8,
				"w": 12,
				"x": (i % 2) * 12,
				"y": (i / 2) * 8,
			},
			"targets": []map[string]any{
				{
					"refId": "A",
					// 这里的表达式需要根据实际的数据源和查询逻辑来确定
					// 暂时使用一个 placeholder
					"expr": fmt.Sprintf(`iot_data{topic="%s", field="%s"}`, unsDef.Topic, field.Name),
				},
			},
		}
		panels = append(panels, panel)
	}

	dashboard := map[string]any{
		"dashboard": map[string]any{
			"uid":         uid,
			"title":       unsDef.Name,
			"description": "Auto-generated by supos-edge for UNS: " + unsDef.Alias,
			"panels":      panels,
			"time": map[string]string{
				"from": "now-6h",
				"to":   "now",
			},
			"timezone": "browser",
		},
		"folderUid": "", // 可以指定一个 folder
		"overwrite": false,
	}

	jsonBytes, err := json.Marshal(dashboard)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// DataExport 导出 Dashboard 数据
func (s *DashboardService) DataExport(exportParam *dto.DashboardExportParam) (string, error) {
	context := &exporter.DashboardExportContext{}

	// 1. 获取数据
	err := s.fetchDataForExport(context, exportParam)
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
func (s *DashboardService) fetchDataForExport(context *exporter.DashboardExportContext, exportParam *dto.DashboardExportParam) error {
	var dashboards []*model.DashboardModel
	var err error

	if len(exportParam.Ids) > 0 {
		dashboards, err = s.dashboardMapper.SelectByIds(exportParam.Ids)
	} else if exportParam.ExportType == "ALL" {
		dashboards, err = s.dashboardMapper.SelectAll()
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
func (s *DashboardService) AsyncImport(filePath string) (*dto.RunningStatus, error) {
	file, err := os.Open(filePath)
	if err != nil {
		s.logger.Errorf("failed to open import file %s: %v", filePath, err)
		return dto.NewRunningStatus(400, "global.import.file.not.exist"), nil
	}
	defer file.Close()

	context := importer.NewDashboardImportContext(filePath)
	dataImporter := importer.NewDashboardDataImporter(context, s.dashboardMapper)

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

	if context.DataEmpty() {
		return dto.NewRunningStatus(400, "dashboard.import.excel.empty"), nil
	}

	if len(context.CheckErrorMap) == 0 {
		status := dto.NewRunningStatus(200, "dashboard.import.rs.ok").SetTask(finalTask).SetProgress(100.0)
		status.TotalCount = context.Total
		status.SuccessCount = context.Total
		status.ErrorCount = 0
		return status, nil
	}

	// 存在部分错误
	relativePath, err := s.writeImportErrorFile(dataImporter)
	if err != nil {
		s.logger.Errorf("failed to write import error file: %v", err)
		// 即使写入错误文件失败，也要通知前端导入已完成但有错误
		status := dto.NewRunningStatus(206, "dashboard.import.rs.hasErr").SetTask(finalTask).SetProgress(100.0)
		status.TotalCount = context.Total
		status.ErrorCount = len(context.CheckErrorMap)
		status.SuccessCount = status.TotalCount - status.ErrorCount
		return status, nil
	}

	message := "dashboard.import.rs.hasErr"
	status := dto.NewRunningStatus(206, message, relativePath).SetTask(finalTask).SetProgress(100.0)
	status.TotalCount = context.Total
	status.ErrorCount = len(context.CheckErrorMap)
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

// 辅助函数

// camelToSnake 驼峰转下划线
func camelToSnake(s string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	return strings.ToLower(re.ReplaceAllString(s, "${1}_${2}"))
}
