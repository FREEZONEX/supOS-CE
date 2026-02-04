package sourceflow

import (
	dao "backend/internal/repo/relationDB"
	"bytes"
	"testing"

	"gitee.com/unitedrhino/share/conf"
)

func TestNodeRedFlowExport(t *testing.T) {
	dao.InitDbConfig(conf.Database{DSN: "postgres://postgres:postgres@100.100.100.20:31014/postgres?search_path=supos"})
	defer removeMock()
	initMock(t)

	groupIds := []int64{3, 4}
	ids := []int64{2004190137830871040, 2004190506053013504, 1963148065022947328, 1910955436975697920}
	{
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		NodeRedFlowExport(t.Context(), groupIds, ids, true, "http://nodered:1880")(buf)
		t.Log(buf.String())
	}
}
