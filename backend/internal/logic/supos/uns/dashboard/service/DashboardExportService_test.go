package service

import (
	"backend/internal/common"
	"backend/internal/config"
	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"gitee.com/unitedrhino/share/conf"
	"github.com/h2non/gock"
	"github.com/zeromicro/go-zero/core/logx"
)

const GRAFANA_HOST = "http://grafana:3000"

func TestNewDashboardExport(t *testing.T) {
	dao.InitDbConfig(conf.Database{DSN: "postgres://postgres:postgres@100.100.100.20:31014/postgres?search_path=supos"})
	spring.RegisterBean(&svc.ServiceContext{
		Config: config.Config{
			GrafanaUrl: GRAFANA_HOST,
		},
	})
	groupIds := []int64{3, 4}
	ids := []string{"09562aa939c944ee", "9e6973984ffb4b8d"}
	req := &types.GlobalExportParam{
		DashboardExportParam: &types.DashboardExportParam{GroupIds: groupIds, DashIds: ids},
	}
	{
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		serv := &DashboardExportService{}
		serv.ExportStream(t.Context(), req)(buf)
		t.Log(buf.String())
	}
}

func TestNewDashboardImport(t *testing.T) {
	defer removeMock()
	initMock(t)
	spring.RegisterBean(&svc.ServiceContext{
		Config: config.Config{
			GrafanaUrl: GRAFANA_HOST,
		},
	})
	jsonCotent := `{
    "data": [
{"id":"85e9e6aed5485541","name":"hktest/dev/Metric/风物","jsonContent":"{\"meta\":{\"type\":\"db\",\"canSave\":true,\"canEdit\":true,\"canAdmin\":true,\"canStar\":false,\"canDelete\":true,\"slug\":\"hktest-dev-metric-e9a38e-e789a9\",\"url\":\"/grafana/home/d/85e9e6aed5485541/hktest-dev-metric-e9a38e-e789a9\",\"expires\":\"0001-01-01T00:00:00Z\",\"created\":\"2026-03-05T03:19:21Z\",\"updated\":\"2026-03-05T03:19:21Z\",\"updatedBy\":\"Anonymous\",\"createdBy\":\"Anonymous\",\"version\":1,\"hasAcl\":false,\"isFolder\":false,\"folderId\":0,\"folderUid\":\"\",\"folderTitle\":\"General\",\"folderUrl\":\"\",\"provisioned\":false,\"provisionedExternalId\":\"\",\"annotationsPermissions\":{\"dashboard\":{\"canAdd\":true,\"canEdit\":true,\"canDelete\":true},\"organization\":{\"canAdd\":true,\"canEdit\":true,\"canDelete\":true}}},\"dashboard\":{\"annotations\":{\"list\":[{\"builtIn\":1,\"datasource\":{\"type\":\"grafana\",\"uid\":\"-- Grafana --\"},\"enable\":true,\"hide\":true,\"iconColor\":\"rgba(0, 211, 255, 1)\",\"name\":\"Annotations \\u0026 Alerts\",\"type\":\"dashboard\"}]},\"editable\":true,\"fiscalYearStartMonth\":0,\"graphTooltip\":0,\"id\":4,\"links\":[],\"panels\":[{\"datasource\":{\"type\":\"grafana-postgresql-datasource\",\"uid\":\"93dec6b2190b0d07\"},\"fieldConfig\":{\"defaults\":{\"color\":{\"mode\":\"palette-classic\"},\"custom\":{\"axisBorderShow\":false,\"axisCenteredZero\":false,\"axisColorMode\":\"text\",\"axisLabel\":\"\",\"axisPlacement\":\"auto\",\"barAlignment\":0,\"barWidthFactor\":0.6,\"drawStyle\":\"line\",\"fillOpacity\":0,\"gradientMode\":\"none\",\"hideFrom\":{\"legend\":false,\"tooltip\":false,\"viz\":false},\"insertNulls\":false,\"lineInterpolation\":\"linear\",\"lineWidth\":1,\"pointSize\":5,\"scaleDistribution\":{\"type\":\"linear\"},\"showPoints\":\"auto\",\"spanNulls\":false,\"stacking\":{\"group\":\"A\",\"mode\":\"none\"},\"thresholdsStyle\":{\"mode\":\"off\"}},\"mappings\":[],\"thresholds\":{\"mode\":\"absolute\",\"steps\":[{\"color\":\"green\"},{\"color\":\"red\",\"value\":80}]}},\"overrides\":[]},\"gridPos\":{\"h\":18,\"w\":24,\"x\":0,\"y\":0},\"id\":1,\"options\":{\"legend\":{\"calcs\":[],\"displayMode\":\"list\",\"placement\":\"bottom\",\"showLegend\":true},\"tooltip\":{\"mode\":\"single\",\"sort\":\"none\"}},\"pluginVersion\":\"11.4.0\",\"targets\":[{\"datasource\":{\"type\":\"grafana-postgresql-datasource\",\"uid\":\"93dec6b2190b0d07\"},\"editorMode\":\"code\",\"format\":\"table\",\"rawQuery\":true,\"rawSql\":\"SELECT * FROM \\\"public\\\".\\\"_fengwu_4d4f3dbf8ea4461b82f0\\\" where 1=1  and $__timeFilter(\\\"timeStamp\\\") \",\"refId\":\"A\",\"sql\":{\"columns\":[{\"parameters\":[],\"type\":\"function\"}],\"groupBy\":[{\"property\":{\"type\":\"string\"},\"type\":\"groupBy\"}],\"limit\":50},\"table\":\"_gg_393282ae8afd4b0c9ebe\"},{\"datasource\":{\"type\":\"grafana-postgresql-datasource\",\"uid\":\"93dec6b2190b0d07\"},\"hide\":false,\"refId\":\"B\"}],\"title\":\"hktest/dev/Metric/风物\",\"type\":\"timeseries\"}],\"preload\":false,\"schemaVersion\":40,\"tags\":[],\"templating\":{\"list\":[]},\"time\":{\"from\":\"now-5m\",\"to\":\"now\"},\"timepicker\":{},\"timezone\":\"browser\",\"title\":\"hktest/dev/Metric/风物\",\"uid\":\"85e9e6aed5485541\",\"version\":1,\"weekStart\":\"\"}}","creator":"tier0","updateTime":"2026-03-05T03:19:21.729179Z","createTime":"2026-03-05T03:19:21.729015Z"}
    ]
}`
	saveGroup := func(ctx context.Context, list []*dao.GroupModel) error {
		bs, _ := json.Marshal(list)
		t.Logf("数据库 saveGroups: %v", string(bs))
		return nil
	}
	dashSave := func(ctx context.Context, list []*dao.DashboardModel) error {
		bs, _ := json.Marshal(list)
		t.Logf("数据库 saveDashboards: %v", string(bs))
		return nil
	}
	progress := common.Float3(0)
	decoder := json.NewDecoder(bytes.NewBufferString(jsonCotent))
	_, err := decoder.Token()
	if err != nil {
		t.Error(err)
	}
	for decoder.More() {
		fieldName, _ := decoder.Token()
		propName, isString := fieldName.(string)
		if !isString {
			logx.WithContext(t.Context()).Errorf("未知Token :%v", fieldName)
			// 跳过未知字段的值
			continue
		}
		switch propName {
		case "data":
			importDashboards(t.Context(), int64(len(jsonCotent)), &progress, func(status *common.RunningStatus) {
				bs, _ := json.Marshal(status)
				t.Log("进度：", string(bs))
			}, decoder, saveGroup, dashSave)
		}
	}

}
func removeMock() {
	gock.Off()
}
func initMock(t *testing.T) {
	gock.New(GRAFANA_HOST).
		Post("/api/dashboards/db").Reply(200).JSON(map[string]string{"code": "200", "msg": "保存成功"})
}
