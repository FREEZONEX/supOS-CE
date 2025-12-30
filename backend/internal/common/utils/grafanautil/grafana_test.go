package grafanautil

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGetDashboardUUIDByAlias(t *testing.T) {
	uuid := GetDashboardUUIDByAlias("_m2_f0afdf2d1b0b45ac820a")
	t.Log(uuid)
}
func TestFormatTemplateMap(t *testing.T) {
	template := LoadTemplate(t.Context(), "templates/td-dashboard-list.json")
	uid := uuid.New().String()[:32] // Fast simple UUID

	panelJSONList := make([]any, 0, 10)
	i := 0
	gridPosX := i * 8
	if gridPosX > 16 {
		gridPosX = (i % 3) * 8
	}
	panelParam := map[string]any{
		"id":        1,
		"title":     "标题",
		"columns":   "columns",
		"schema":    "public",
		"tableName": 11111,
		"gridPosX":  gridPosX,
	}
	panelTemplate := LoadTemplate(t.Context(), "templates/ts-panel.json")

	panelTemplate = strings.Replace(panelTemplate, "supos_tag_timeserial", "TableX", 1)
	panelTemplate = strings.Replace(panelTemplate, "tag_name", "tag", 1)
	panelJSON := FormatTemplateMap(panelTemplate, panelParam)
	var panel any
	json.Unmarshal([]byte(panelJSON), &panel)
	panelJSONList = append(panelJSONList, panel)

	panelJSONList = append(panelJSONList, panel)

	dbParams := map[string]any{
		"uid":    uid,
		"title":  "整个Title",
		"panels": panelJSONList,
	}

	dashboardJSON := FormatTemplateMap(template, dbParams)
	t.Logf("创建时序组合DashboardDashboard 请求: %s", dashboardJSON)
}

func TestCreateDs(t *testing.T) {
	t.Log(GetDashboardUUIDByAlias("Metric/周六晚9点"))
	t.Log(GetDashboardUUIDByAlias("_zhouliuwan9dian_b33eec219ea4431e914a"))

	//t.Log(GetDataSourceByName(types.SrcJdbcTypeTimeScaleDB.Alias()))
}
