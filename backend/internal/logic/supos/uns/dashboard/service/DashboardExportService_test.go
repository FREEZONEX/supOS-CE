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
        {
            "id": "09562aa939c944ee",
            "name": "_zz03",
            "type": 1,
            "updateTime": "2025-05-15T02:49:11.057Z",
            "createTime": "2025-05-15T02:49:11.057Z",
            "jsonContent": "{}"
        },
        {
            "id": "1",
            "name": "分组1",
            "exportType": "group",
            "description": "硕鼠所",
            "updateTime": "2026-01-21T08:56:54.722Z",
            "createTime": "2026-01-21T08:56:54.722Z",
            "children": [
                {
                    "id": "9e6973984ffb4b8d",
                    "name": "Metric/劳平二提量",
                    "creator": "guest",
                    "updateTime": "2025-12-27T12:36:19.199776Z",
                    "createTime": "2025-12-27T12:36:19.199776Z"
                }
            ]
        }
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
	importDashboards(t.Context(), int64(len(jsonCotent)), func(status *common.RunningStatus) {
		bs, _ := json.Marshal(status)
		t.Log("进度：", string(bs))
	}, bytes.NewBufferString(jsonCotent), saveGroup, dashSave)
}
func removeMock() {
	gock.Off()
}
func initMock(t *testing.T) {
	gock.New(GRAFANA_HOST).
		Post("/api/dashboards/db").Reply(200).JSON(map[string]string{"code": "200", "msg": "保存成功"})
}
