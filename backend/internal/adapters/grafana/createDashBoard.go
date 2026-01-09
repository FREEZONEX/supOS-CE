package grafana

import (
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"fmt"
	"strings"
	"time"
)

// Create 创建仪表板的主方法
func (g *GrafanaEventHandler) Create(ctx context.Context, jdbcType types.SrcJdbcType, topics []*types.CreateTopicDto, flowName string, fromImport bool, username string) {
	// 和 CE 一样，现在不区分导入的情况

	storageAdapter := g.getPersistentService(jdbcType)
	if storageAdapter == nil {
		g.log.Errorf("忽略不支持的 jdbcType=%s", jdbcType.Alias)
		return
	}

	ds := storageAdapter.GetDataSourceProperties()
	startTime := time.Now()
	g.log.Infof("开始创建仪表板，预计数量：%d, 流程：%s, 数据源类型：%s", len(topics), flowName, jdbcType.Alias)

	//var seqMergeUns []*types.CreateTopicDto
	var normalUns = make([]*types.CreateTopicDto, 0, len(topics))
	//isSeq := jdbcType.TypeCode() == constants.TimeSequenceType
	for _, topic := range topics {
		if !base.P2v(topic.AddDashBoard) {
			continue
		}
		/*	if fromImport && isSeq && topic.GetTbFieldName() != "" {
				seqMergeUns = append(seqMergeUns, topic)
			} else {
			}*/
		normalUns = append(normalUns, topic)
	}
	if len(normalUns) > 0 {

		createdCount := g.createIndividualDashboards(ctx, normalUns, ds, jdbcType, fromImport, username)

		elapsed := time.Since(startTime)
		g.log.Infof("完成创建仪表板，实际创建：%d, 耗时：%v, 流程：%s", createdCount, elapsed, flowName)
	}

	//// 处理组合仪表板
	//if fromImport && len(seqMergeUns) > 0 {
	//	g.createCompositeDashboard(ctx, jdbcType, seqMergeUns, flowName, username)
	//}
}

// createIndividualDashboards 创建单个仪表板
func (g *GrafanaEventHandler) createIndividualDashboards(ctx context.Context, topics []*types.CreateTopicDto, ds serviceApi.DataSourceProperties, jdbcType types.SrcJdbcType, fromImport bool, username string) int {

	dashboards := make([]event.DashboardVo, 0, len(topics))
	for _, dto := range topics {
		if !base.P2v(dto.AddDashBoard) {
			continue
		}
		columns := grafanautil.Fields2Columns(jdbcType, dto.Fields)
		title := dto.Path
		schema, table := g.extractSchemaAndTable(dto.Alias, ds.Schema)
		tagNameCondition := ""

		g.log.Debugf("创建 Grafana 仪表板 - 列：%s, 标题：%s, 模式：%s, 表：%s, 标签条件：%s, fromImport? %v",
			columns, title, schema, table, tagNameCondition, fromImport)

		dashId := grafanautil.GetDashboardUUIDByAlias(dto.Alias)
		err := grafanautil.CreateDashboard(dashId, ctx, table, tagNameCondition, jdbcType, schema, title, columns, constants.SysFieldCreateTime)
		if err != nil {
			g.log.Infof("创建仪表板失败: %v", err)
			continue
		}
		//  发布事件
		desc := base.P2v(dto.Description)
		if len(desc) == 0 {
			desc = dto.Path
		}
		dashboards = append(dashboards, event.DashboardVo{
			UUID:        dashId,
			UnsAlias:    []string{dto.Alias},
			Name:        title,
			Description: desc,
			UserName:    username,
		})
	}
	if len(dashboards) > 0 {
		er := spring.PublishEvent(&event.CreateDashboardEvent{
			ApplicationEvent: event.ApplicationEvent{Context: ctx},
			Dashboards:       dashboards,
		})
		if er != nil {
			g.log.Error("CreateDashboardErr", er)
		}
	}
	return len(dashboards)
}

// extractSchemaAndTable 从表名中提取模式和表名
func (g *GrafanaEventHandler) extractSchemaAndTable(fullTableName, defaultSchema string) (string, string) {
	if dot := strings.Index(fullTableName, "."); dot > 0 {
		return fullTableName[:dot], fullTableName[dot+1:]
	}
	return defaultSchema, fullTableName
}

// buildTagNameCondition 构建标签名称条件
//func (g *GrafanaEventHandler) buildTagNameCondition(dto *types.CreateTopicDto) string {
//	if dto.GetTbFieldName() == "" {
//		return ""
//	}
//	return fmt.Sprintf(" and %s='%d' ", constants.SystemSeqTag, dto.Id)
//}

// createCompositeDashboard 创建组合仪表板
func (g *GrafanaEventHandler) createCompositeDashboard(ctx context.Context, jdbcType types.SrcJdbcType, topics []*types.CreateTopicDto, flowName, username string) {
	if jdbcType.TypeCode() != constants.TimeSequenceType {
		return
	}

	// 获取临时名称
	tempName := g.getTempName(topics)
	if tempName == "" {
		return
	}

	dashboardName := generateDashboardName(flowName, tempName)
	uuid, err := grafanautil.CreateTimeSeriesListDashboard(ctx, jdbcType, topics, dashboardName)
	if err != nil {
		g.log.Error("创建时序面板失败", err.Error())
		return
	}
	unsAliasList := base.Map[*types.CreateTopicDto, string](topics, func(e *types.CreateTopicDto) string {
		return e.Alias
	})
	spring.PublishEvent(event.NewCreateDashboardEvent(ctx, unsAliasList, uuid, dashboardName, "组合仪表板", username))
}

// getTempName 获取临时名称
func (g *GrafanaEventHandler) getTempName(topics []*types.CreateTopicDto) string {
	for _, dto := range topics {
		if base.P2v(dto.AddDashBoard) {
			return dto.Alias
		}
	}
	return ""
}

// generateDashboardName 生成仪表板名称
func generateDashboardName(flowName, tempName string) string {
	timestamp := time.Now().Format("20060102150405")
	if flowName != "" {
		return fmt.Sprintf("%s-%s", flowName, timestamp)
	}
	return fmt.Sprintf("%s-%s", tempName, timestamp)
}
