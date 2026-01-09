package dashboard

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/utils/grafanautil"
	unsservice "backend/internal/logic/supos/uns/uns/service"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGrafanaByUnsLogic struct {
	logx.Logger
	ctx                  context.Context
	svcCtx               *svc.ServiceContext
	dashboardMapper      relationDB.DashboardMapper
	dashboardRefMapper   relationDB.DashboardRefMapper
	unsDefinitionService serviceApi.IUnsDefinitionService
	unsAddService        *unsservice.UnsAddService

	unsQueryService  *unsservice.UnsQueryService
	unsUpdateService *unsservice.UnsUpdateService
}

func NewCreateGrafanaByUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGrafanaByUnsLogic {
	// Note: Services might be managed by spring, adjust if needed

	return &CreateGrafanaByUnsLogic{
		Logger:               logx.WithContext(ctx),
		ctx:                  ctx,
		svcCtx:               svcCtx,
		unsDefinitionService: spring.GetBean[serviceApi.IUnsDefinitionService](),
		unsAddService:        spring.GetBean[*unsservice.UnsAddService](),

		unsQueryService:  spring.GetBean[*unsservice.UnsQueryService](),
		unsUpdateService: spring.GetBean[*unsservice.UnsUpdateService](),
	}
}
func (l *CreateGrafanaByUnsLogic) CreateGrafanaByUns(alias string) (*types.JsonResult, error) {
	uns := l.unsDefinitionService.GetDefinitionByAlias(alias)
	if uns == nil {
		return &types.JsonResult{Code: 400, Msg: I18nUtils.GetMessage("uns.file.not.exist")}, nil
	}
	jdbcType := types.SrcJdbcType(uns.DataSrcID)
	columns := grafanautil.Fields2Columns(jdbcType, uns.Fields)
	title := uns.Path
	schema := "public"
	table := uns.Alias
	tagNameCondition := ""
	// 日志输出（Go中使用log包）
	l.Logger.Debugf(">>>>> create grafana dashboard columns:%s,title:%s,schema:%s,table:%s,tagNameCondition:%s",
		columns, title, schema, table, tagNameCondition)
	dot := strings.Index(table, ".")
	if dot > 0 {
		schema = table[:dot]
		table = table[dot+1:]
	}
	dashId := grafanautil.GetDashboardUUIDByAlias(alias)
	err := grafanautil.CreateDashboard(dashId, l.ctx, table, tagNameCondition, jdbcType, schema, title, columns, constants.SysFieldCreateTime)
	if err != nil {
		return &types.JsonResult{Code: 400, Msg: err.Error()}, nil
	}
	// 检查仪表板是否已存在
	db := relationDB.GetDb(l.ctx)
	dashboardPo, _ := l.dashboardMapper.SelectById(db, dashId)
	if dashboardPo == nil {
		now := time.Now()
		po := &relationDB.DashboardModel{
			ID:         dashId,
			Name:       alias,
			CreateTime: now,
			UpdateTime: now,
		}

		// 设置创建者
		if userCtx := apiutil.GetUserFromContext(l.ctx); userCtx != nil {
			po.Creator = userCtx.PreferredUsername
		}

		if err := l.dashboardMapper.Insert(db, po); err != nil {
			fmt.Printf("Failed to insert dashboard: %v\n", err)
		}
	}

	// 生成标志位并创建分类模型实例
	if uns.WithFlags == nil {
		var flag int32 = constants.UnsFlagWithDashboard
		uns.WithFlags = &flag
	} else {
		*uns.WithFlags |= constants.UnsFlagWithDashboard
	}

	if err := l.unsAddService.CreateModelInstance(l.ctx, uns); err != nil {
		fmt.Printf("Failed to create category model instance: %v\n", err)
	}

	// 创建仪表板引用关系（如果不存在）
	refPo, _ := l.dashboardRefMapper.GetByUns(db, alias)
	if refPo == nil {
		ref := &relationDB.DashboardRefModel{
			DashboardID: dashId,
			UnsAlias:    alias,
		}
		if err := l.dashboardRefMapper.Insert(db, ref); err != nil {
			fmt.Printf("Failed to insert dashboard reference: %v\n", err)
		}
	}

	return &types.JsonResult{
		Code: 200,
		Msg:  "ok",
		Data: dashId,
	}, nil
}
func (l *CreateGrafanaByUnsLogic) CreateGrafanaByUns_old(alias string) (*types.JsonResult, error) {
	// 1. 获取 UNS 定义
	unsResp, err := l.unsQueryService.GetModelDefinition(l.ctx, &types.ModelDetailReq{}, alias)
	if err != nil || unsResp == nil || unsResp.Data == nil {
		l.Logger.Errorf("failed to get uns definition for alias %s: %v", alias, err)
		return &types.JsonResult{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}
	unsDef := unsResp.Data

	db := relationDB.GetDb(l.ctx)
	// 2. 检查是否已有关联的 Dashboard
	existingDashboard, err := l.dashboardRefMapper.GetByUns(db, alias)
	if err != nil {
		l.Logger.Errorf("error checking for existing dashboard for uns %s: %v", alias, err)
		// Fall through, but log the error
	}
	if existingDashboard != nil {
		return &types.JsonResult{
			Code: http.StatusBadRequest,
			Msg:  i18ns.LocalizeMsg("uns.dashboard.already.exists"),
		}, nil
	}

	// 3. 根据 UNS 字段构建 Grafana Dashboard JSON
	dashboardUID := uuid.New().String()
	dashboardJSON, err := l.buildGrafanaJSONFromUns(unsDef, dashboardUID)
	if err != nil {
		l.Logger.Errorf("failed to build grafana json for uns %s: %v", alias, err)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}

	// 4. 调用 Grafana API 创建 Dashboard
	_, err = grafanautil.CreateDashboardByBody(l.ctx, dashboardUID, "", dashboardJSON)
	if err != nil {
		l.Logger.Errorf("failed to create grafana dashboard for uns %s: %v", alias, err)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}

	// 5. 保存 Dashboard 记录到数据库
	now := time.Now()
	dashboard := &relationDB.DashboardModel{
		ID:          dashboardUID,
		Name:        unsDef.Name,
		Creator:     "system", // 系统自动创建
		CreateTime:  now,
		UpdateTime:  now,
		Type:        1, // 1 for Grafana
		JsonContent: dashboardJSON,
		NeedInit:    false, // 已在 Grafana 中创建
	}
	if err = l.dashboardMapper.Insert(db, dashboard); err != nil {
		l.Logger.Errorf("failed to save dashboard record for uns %s: %v", alias, err)
		// 尝试回滚 Grafana 的创建操作
		_ = grafanautil.DeleteDashboard(l.ctx, dashboardUID)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}

	// 6. 创建 Dashboard 和 UNS 的绑定关系
	ref := &relationDB.DashboardRefModel{
		DashboardID: dashboardUID,
		UnsAlias:    alias,
		CreateAt:    now,
	}
	if err = l.dashboardRefMapper.Insert(db, ref); err != nil {
		l.Logger.Errorf("failed to bind dashboard %s to uns %s: %v", dashboardUID, alias, err)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}, nil
	}

	// 7. 更新 UNS 的 flags
	updateDto := &types.UpdateUnsDto{
		Alias:        alias,
		AddDashBoard: base.OptionalTrue,
	}
	_, updateErr := l.unsUpdateService.UpdateDetail(l.ctx, updateDto)
	if updateErr != nil {
		l.Logger.Errorf("failed to update uns flags for %s after creating dashboard: %v", alias, updateErr)
		return &types.JsonResult{
			Code: http.StatusInternalServerError,
			Msg:  updateErr.Error(),
		}, nil
	}

	// 8. 返回新创建的 Dashboard UID
	return &types.JsonResult{
		Code: http.StatusOK,
		Msg:  "success",
		Data: dashboardUID,
	}, nil
}

// buildGrafanaJSONFromUns 根据 UNS 定义构建 Grafana Dashboard JSON
func (l *CreateGrafanaByUnsLogic) buildGrafanaJSONFromUns(unsDef *types.ModelDetail, uid string) (string, error) {
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
					"expr":  fmt.Sprintf(`iot_data{topic="%s", field="%s"}`, unsDef.Topic, field.Name),
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
		"folderUid": "",
		"overwrite": false,
	}

	jsonBytes, err := json.Marshal(dashboard)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
