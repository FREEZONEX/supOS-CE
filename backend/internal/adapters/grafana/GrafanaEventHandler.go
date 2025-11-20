package grafana

import (
	sysconfig "backend/internal/common/config"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrafanaEventHandler struct {
	log       logx.Logger
	once      sync.Once
	sysConfig *sysconfig.SystemConfig
	dsMap     map[types.SrcJdbcType]serviceApi.IPersistentService
}

func init() {
	spring.RegisterLazy[*GrafanaEventHandler](func() *GrafanaEventHandler {
		sysConfig := spring.GetBean[*sysconfig.SystemConfig]()
		log := logx.WithContext(context.Background())
		if !base.MapContainsKey(sysConfig.ContainerMap, "grafana") {
			log.Debug(">>>>>>>>>当前系统未启用grafana服务!")
			return nil
		}
		return &GrafanaEventHandler{
			log:       log,
			sysConfig: sysConfig,
		}
	})
}
func (g *GrafanaEventHandler) getPersistentService(dsId types.SrcJdbcType) serviceApi.IPersistentService {
	if g.dsMap == nil {
		g.once.Do(func() {
			g.dsMap = base.MapArrayToMap(spring.GetBeansOfType[serviceApi.IPersistentService](),
				func(e serviceApi.IPersistentService) (ok bool, k types.SrcJdbcType, v serviceApi.IPersistentService) {
					return true, e.GetDataSrcId(), e
				})
		})
	}
	return g.dsMap[dsId]
}
func (g *GrafanaEventHandler) OnEventBatchCreateTable300(evt *event.BatchCreateTableEvent) {
	g.log.Infof(">>>>>> GrafanaEventHandler 批量创建事件,topic数量：%d,flowName:%s", len(evt.Creates), evt.FlowName)
	if len(evt.Creates) == 0 {
		return
	}
	userCtx := apiutil.GetUserFromContext(evt.Context)
	userName := ""
	if userCtx != nil {
		userName = userCtx.PreferredUsername
	}
	go func() {
		for dsId, list := range evt.GetAllCreateFiles() {
			g.Create(evt.Context, dsId, list, evt.FlowName, evt.FromImport, userName)
		}
		g.log.Info(">>>>>> GrafanaEventHandler 批量创建事件,已完成,flowName:", evt.FlowName)
	}()
}
func (g *GrafanaEventHandler) OnEventUpdateInstanceEvent300(evt *event.UpdateInstanceEvent) {
}
func (g *GrafanaEventHandler) OnEventRemoveTopicsEvent300(evt *event.RemoveTopicsEvent) {
	list := evt.Topics
	if len(list) == 0 || !evt.WithDashboard {
		return
	}
	tables := base.FilterAndMap[*types.CreateTopicDto, string](list, func(e *types.CreateTopicDto) (v string, ok bool) {
		return e.GetTable(), e.PathType == constants.PathTypeFile && !constants.WithRetainTableWhenDeleteInstance(base.P2v(e.WithFlags))
	})
	go func() {
		for _, table := range tables {
			uid := grafanautil.GetDashboardUUIDByAlias(table)
			err := grafanautil.DeleteDashboard(uid)
			if err != nil {
				g.log.Error("删除grafana dashboard 异常", err, table)
			}
		}
	}()
}
func (g *GrafanaEventHandler) OnEventContextRefreshedEvent300(evt *event.ContextRefreshedEvent) {
	stop := false
	go func() {
		time.Sleep(time.Second)
		g.getPersistentService(0)
		for !stop {
			countOk := 0
			for srcId, ds := range g.dsMap {
				rsCode := grafanautil.GetDataSourceByName(srcId.DataSrcType())
				if rsCode == 200 {
					countOk++
				} else {
					ok, err := grafanautil.CreateDatasource(srcId, ds.GetDataSourceProperties(), false)
					if err != nil {
						g.log.Error("CreateDatasourceErr:", srcId.Alias(), err)
					} else if ok {
						countOk++
					}
				}
			}
			if countOk == len(g.dsMap) {
				break
			}
			time.Sleep(15 * time.Second)
		}
		if "zh-CN" == g.sysConfig.Lang {
			_ = grafanautil.SetLanguage("zh-Hans")
		}
		g.log.Info(">>>>>>>>>Grafana 默认datasource 完成创建.")
	}()

}
