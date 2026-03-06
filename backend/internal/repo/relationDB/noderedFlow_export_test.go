package relationDB

import (
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	"backend/share/base"
	"bytes"
	"fmt"
	"io"
	"testing"
)

func init() {
	dbConfig.DSN = "postgres://postgres:postgres@100.100.100.20:31014/postgres?search_path=supos"
}
func TestExportEmpty(t *testing.T) {
	var flowUnsMapper NodeRedFlowUnsExporter
	flowUnsCsv2po := func(headers, values []string) *NoderedFlowNode {
		return flowUnsMapper.Csv2Model(headers, values)
	}
	flowUnsGetChildren := func(node *NoderedFlowNode) []*NoderedFlowNode {
		return nil
	}
	flowUnsSetChildren := func(node *NoderedFlowNode, children []*NoderedFlowNode) {
	}
	flowUnsGetId := func(node *NoderedFlowNode) int64 {
		return node.ParentID
	}
	flowUnsGetParentId := func(node *NoderedFlowNode) int64 {
		return -1
	}
	jsonWriter := bytes.NewBuffer([]byte{})
	_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
		return flowUnsMapper.ExportByFlowIds(t.Context(), []int64{1, 3}, writer)
	}, jsonWriter, flowUnsGetChildren, flowUnsSetChildren, flowUnsGetId, flowUnsGetParentId, flowUnsCsv2po, true)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(jsonWriter.String())
	}
}
func Test_ExportFlowByGroupAndIds(t *testing.T) {
	groupIds := []int64{3, 4}
	ids := []int64{2004190137830871040, 2004190506053013504, 2015680305511272448, 2015680315900563456}
	logFlowJson(t, groupIds, ids)
	logFlowJson(t, groupIds, nil)
	logFlowJson(t, nil, ids)
	logFlowJson(t, nil, nil)
}
func logFlowJson(t *testing.T, groupIds []int64, ids []int64) {
	if len(groupIds) != 0 && len(ids) != 0 {
		t.Log("export 分组和id: ")
	} else if len(groupIds) != 0 {
		t.Log("export 仅分组: ")
	} else if len(ids) != 0 {
		t.Log("export 仅id: ")
	} else {
		t.Log("导出全部")
	}
	mapper := NodeRedFlowExporter{}
	csv2po := func(headers, values []string) *NoderedFlow {
		return mapper.Csv2Model(headers, values)
	}
	jsonWriter := bytes.NewBuffer([]byte{})
	fmt.Fprintln(jsonWriter, `{ "data":`)
	_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
		return mapper.ExportByGroupAndIds(t.Context(), groupIds, ids, true, writer)
	}, jsonWriter, flowGetChildren, flowSetChildren, flowGetId, flowGetParentId, csv2po, true)
	if err == nil {
		fmt.Fprintln(jsonWriter, `}`)
	} else {
		t.Error("Csv2JsonStream err:", err)
	}
	t.Log(jsonWriter.String())

}
func flowGetChildren(node *NoderedFlow) []*NoderedFlow {
	return node.Children
}
func flowSetChildren(node *NoderedFlow, children []*NoderedFlow) {
	node.Children = children
}
func flowGetId(node *NoderedFlow) int64 {
	return node.ID
}
func flowGetParentId(node *NoderedFlow) int64 {
	return base.P2v(node.GroupId)
}
