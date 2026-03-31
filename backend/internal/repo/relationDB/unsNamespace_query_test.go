package relationDB

import (
	"backend/internal/common/dto"
	"backend/internal/config"
	"backend/share/base"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
)

func TestJsonErr(t *testing.T) {
	fs := Fields{
		{Name: "timeStamp", Type: "DATETIME", Unique: base.OptionalTrue, SystemField: base.OptionalTrue},
		{Name: "wet", Type: "FLOAT"},
		{Name: "wq", Type: "LONG"},
		{Name: "status", Type: "LONG", SystemField: base.OptionalTrue},
	}
	bs, er := json.Marshal(fs)
	if er != nil {
		t.Fatal(er)
	}
	t.Log(string(bs))
	t.Log(fmt.Sprintf("%v", string(bs)))
}
func TestListByLayRecs(t *testing.T) {
	t.Skip("legacy test referenced ExistsTimeSeriaNoneTables, which is no longer present after upstream merge")
}
func TestLabelListAll(t *testing.T) {
	db := stores.GetCommonConn(t.Context())
	var lb UnsLabelRepo
	lbs, err := lb.ListAll(db, 3, 5)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lbs)
}
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
	if os.Getenv("RUN_RELATIONDB_INTEGRATION_TESTS") != "1" {
		return
	}
	c := config.Config{
		Database: conf.Database{
			IsInitTable: true,
			DBType:      "pgsql",
			DSN:         "postgres://postgres:postgres@192.168.236.101:5432/postgres?search_path=supos",
		},
	}

	stores.InitConn(c.Database)
	Migrate(c.Database)
}
