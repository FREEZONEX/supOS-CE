package relationDB

import (
	"backend/internal/common/dto"
	"backend/internal/config"
	"backend/share/base"
	"encoding/json"
	"fmt"
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
	db := stores.GetCommonConn(t.Context())
	dao := NewUnsNamespaceRepo()
	page := &stores.PageInfo{Page: 1, Size: 10, Orders: []stores.OrderBy{{Field: "lay_rec"}}}
	layRecs := []string{"1960575789291339779", "1965675474571513856"}
	rs, err := dao.ListByLayRecs(db, layRecs, page)
	if err != nil {
		t.Fatal(err)
	}

	index := 0
	t.Log(index, len(rs), rs)
	for len(rs) > 0 {
		index++
		page.Page++
		rs, err = dao.ListByLayRecs(db, layRecs, page)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(index, len(rs), rs)
	}

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
	c := config.Config{
		Database: conf.Database{
			IsInitTable: true,
			DBType:      "pgsql",
			DSN:         "postgres://postgres:postgres@100.100.100.20:31014/postgres?search_path=supos",
		},
	}

	stores.InitConn(c.Database)
	Migrate(c.Database)
}
