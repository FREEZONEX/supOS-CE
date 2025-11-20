package grafanautil

import (
	"backend/internal/adapters/postgresql"
	"backend/internal/config"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"testing"
)

func TestCreateDs(t *testing.T) {
	spring.RegisterBean[*svc.ServiceContext](&svc.ServiceContext{
		Config: config.Config{
			GrafanaUrl: "http://192.168.235.55:3000",
		},
	})
	//t.Log(GetDataSourceByName(types.SrcJdbcTypeTimeScaleDB.Alias()))

	testUrl := "postgres://postgres:postgres@100.100.100.20:31014/postgres"
	ds := postgresql.ParseDbUrlProperties(testUrl)
	ok, err := CreateDatasource(types.SrcJdbcTypeTimeScaleDB, ds, false)
	t.Log(ok, err)

	//t.Log(GetDataSourceByName(types.SrcJdbcTypeTimeScaleDB.Alias()))
}
