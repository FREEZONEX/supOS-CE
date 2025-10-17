package relationDB

import (
	"backend/internal/config"
	"encoding/json"
	"testing"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
)

func TestUnsQuery(t *testing.T) {
	dao := NewUnsNamespaceRepo(nil)
	ctx := t.Context()

	rs, err := dao.ListTimeSeriesFiles(ctx, "", &stores.PageInfo{Page: 1, Size: 10})
	jbs, _ := json.Marshal(rs)
	t.Log(len(rs), string(jbs), err)

	if len(rs) > 0 {
		rs, err = dao.ListUnsByIds(ctx, []int64{rs[0].ID})
		jbs, _ = json.MarshalIndent(rs, "", " ")
		t.Log(string(jbs), err)
	}
}

func init() {
	c := config.Config{
		Database: conf.Database{
			IsInitTable: true,
			DBType:      "pgsql",
			DSN:         "postgres://postgres:postgres@100.100.100.20:31014/postgres",
		},
		DatabaseSchema: "supos",
	}

	stores.InitConn(c.Database)
	Migrate(c.Database, c.DatabaseSchema)
}
