package grafanautil

import (
	"testing"
)

func TestCreateDs(t *testing.T) {
	t.Log(GetDashboardUUIDByAlias("Metric/周六晚9点"))
	t.Log(GetDashboardUUIDByAlias("_zhouliuwan9dian_b33eec219ea4431e914a"))

	//t.Log(GetDataSourceByName(types.SrcJdbcTypeTimeScaleDB.Alias()))
}
