package relationDB

import (
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	"backend/share/base"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"testing"
)

func TestParseFields(t *testing.T) {
	var fMap = parseTagFields(&DashboardModel{}, "gorm", "json")
	for k, v := range fMap {
		t.Logf("%s = %d\n", k, v)
	}
}
func Test_ExportByGroupAndIds(t *testing.T) {
	dbConfig.DSN = "postgres://postgres:postgres@100.100.100.20:31014/postgres?search_path=supos"
	groupIds := []int64{1, 2, 3}
	ids := []string{"a15f4eee2d0e35bd", "9e6973984ffb4b8d"}
	logJson(t, groupIds, ids)
	logJson(t, groupIds, nil)
	logJson(t, nil, ids)
	logJson(t, nil, nil)
}
func logJson(t *testing.T, groupIds []int64, ids []string) {
	if len(groupIds) != 0 && len(ids) != 0 {
		t.Log("export 分组和id: ")
	} else if len(groupIds) != 0 {
		t.Log("export 仅分组: ")
	} else if len(ids) != 0 {
		t.Log("export 仅id: ")
	} else {
		t.Log("导出全部")
	}
	mapper := DashboardMapper{}
	csv2po := func(headers, values []string) *DashboardModel {
		return mapper.Csv2Model(headers, values)
	}
	jsonWriter := bytes.NewBuffer([]byte{})
	fmt.Fprintln(jsonWriter, `{ "data":`)
	_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
		return mapper.ExportByGroupAndIds(t.Context(), groupIds, ids, writer)
	}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, csv2po, true)
	if err == nil {
		fmt.Fprintln(jsonWriter, `}`)
	} else {
		t.Error("Csv2JsonStream err:", err)
	}
	t.Log(jsonWriter.String())

}
func nodeGetChildren(node *DashboardModel) []*DashboardModel {
	if node.Type == -1 && node.Children == nil { //group
		return []*DashboardModel{}
	}
	return node.Children
}
func nodeSetChildren(node *DashboardModel, children []*DashboardModel) {
	node.Children = children
}
func nodeGetId(node *DashboardModel) string {
	return node.ID
}
func nodeGetParentId(node *DashboardModel) string {
	gid := base.P2v(node.GroupId)
	if gid < 1 {
		return ""
	}
	return strconv.FormatInt(gid, 10)
}
