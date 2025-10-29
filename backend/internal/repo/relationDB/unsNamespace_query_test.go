package relationDB

import (
	"backend/internal/common/dto"
	"backend/internal/config"
	"encoding/json"
	"testing"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
)

func TestUnsQuery(t *testing.T) {
	dao := NewUnsNamespaceRepo()
	db := stores.GetCommonConn(t.Context())
	searchCount := int64(0)
	rs, err := dao.ListTimeSeriesFiles(db, "", &stores.PageInfo{Page: 1, Size: 10}, &searchCount)
	jbs, _ := json.Marshal(rs)
	t.Log(len(rs), string(jbs), err)

	if len(rs) > 0 {
		unsPos, err := dao.ListUnsByIds(db, []int64{1960575789291339779, rs[0].Id})
		jbs, _ = json.MarshalIndent(unsPos, "", " ")
		t.Log(string(jbs), err)
	}
	{
		unsPos, err := dao.ListInTemplate(db, "pride")
		jbs, _ = json.Marshal(unsPos)
		t.Log(len(unsPos), string(jbs), err)
	}
}
func TestListInTemplate(t *testing.T) {
	dao := NewUnsNamespaceRepo()
	db := stores.GetCommonConn(t.Context())
	{
		unsPos, err := dao.ListInTemplate(db, "pride")
		jbs, _ := json.Marshal(unsPos)
		t.Log(len(unsPos), string(jbs), err)
	}
	{
		//count, err := dao.ListAlarmRules(db, "pride")
		//t.Log("countAlarm:", count, err)
	}
}
func TestListByConditions(t *testing.T) {
	dao := NewUnsNamespaceRepo()
	db := stores.GetCommonConn(t.Context())
	{
		unsPos, err := dao.ListByConditions(db, &dto.UnsSearchCondition{
			Keyword:   "pride",
			LabelName: "seq",
		})
		jbs, _ := json.Marshal(unsPos)
		t.Log(len(unsPos), string(jbs), err)
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
