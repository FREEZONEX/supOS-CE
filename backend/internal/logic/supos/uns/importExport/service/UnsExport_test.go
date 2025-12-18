package service

import (
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	dao "backend/internal/repo/relationDB"
	"strings"
	"testing"

	"gitee.com/unitedrhino/share/conf"
)

func TestCsv2JsonStream(t *testing.T) {
	// 初始化数据库连接池
	dao.InitDbConfig(conf.Database{
		DBType: "pgsql",
		DSN:    "postgres://postgres:postgres@100.100.100.20:31014/postgres?search_path=supos",
	})
	jsonWriter := &strings.Builder{}
	var l UnsImportExportService
	countNodes, err := jsonstream.Csv2JsonStream(l.labelMapper.ExportCsv, jsonWriter,
		nodeGetChildren, nodeSetChildren, nodeGetId, func(f *FileData) int64 {
			return -1
		}, l.labelCsv2FileData, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("countNodes:", countNodes)
	t.Log(jsonWriter.String())

}
func Test_getLayAndIdsInner(t *testing.T) {
	layRec, ids := getLayAndIdsInner([]int64{3}, []int64{4, 5}, []string{"1/2/3", "1/2/3/4", "1/2/5"})
	t.Log("layRec = ", layRec)
	t.Log("ids = ", ids)
	//layRec =  [1/2/3]
	//ids =  [5]
}
