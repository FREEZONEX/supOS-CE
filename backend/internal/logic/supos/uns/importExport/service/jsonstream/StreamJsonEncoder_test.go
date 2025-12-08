package jsonstream

import (
	"backend/internal/common/constants"
	"backend/share/base"
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"
)

type IdNode struct {
	ID       int64     `json:"id"`
	ParentId int64     `json:"parentId,omitzero"`
	Name     string    `json:"name"`
	Children []*IdNode `json:"children,omitempty"`

	parent *IdNode `json:"-"`
	Path   string  `json:"path,omitempty"`
}

func (node *IdNode) getPath() string {
	if node.Path == "" {
		if node.parent != nil {
			dir := node.parent.getPath()
			name := node.Name
			if node.Name == "" {
				name = strconv.FormatInt(node.ID, 10)
			}
			node.Path = fmt.Sprintf("%s/%s", dir, name)
		} else {
			name := node.Name
			if node.Name == "" {
				name = strconv.FormatInt(node.ID, 10)
			}
			node.Path = name
		}
	}
	return node.Path
}

func nodeGetId(node *IdNode) int64 {
	return node.ID
}
func nodeGetParentId(node *IdNode) int64 {
	return node.ParentId
}
func nodeGetChildren(node *IdNode) []*IdNode {
	return node.Children
}
func nodeSetChildren(node *IdNode, children []*IdNode) {
	node.Children = children
}
func node2tree(node, parent *IdNode) *TreeNode {
	node.parent = parent
	rs := &TreeNode{ID: node.ID, Name: node.Name, Path: node.getPath()}
	if parent != nil {
		rs.ParentID = parent.ID
	}
	return rs
}

type TreeNode struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	ParentID int64  `json:"parentId"`
}

func (node *TreeNode) String() string {
	if node.Path != "" {
		return node.Path
	}
	if node.ID == 0 && node.Name != "" {
		return fmt.Sprintf("{%s id:%d parent:%d}", node.Name, node.ID, node.ParentID)
	}
	return fmt.Sprintf("{id:%d parent:%d}", node.ID, node.ParentID)
}
func treeGetId(node *TreeNode) int64 {
	return node.ID
}
func treeGetParentId(node *TreeNode) int64 {
	return node.ParentID
}
func tree2node(i *TreeNode) *IdNode {
	return &IdNode{ID: i.ID, Name: i.Name}
}
func Test_getChildrenJsonTagName(t *testing.T) {
	var nt IdNode
	t.Log(getChildrenJsonTagName(nt))
}
func TestStreamDecodeObj(t *testing.T) {
	var nodes = []*IdNode{
		{
			ID: 1,
			Children: []*IdNode{{
				ID: 11,
			}, {ID: 12, Children: []*IdNode{{ID: 120001}, {ID: 120002, Children: []*IdNode{{ID: 12000201}, {ID: 12000202}}}}},

				{ID: 100003, Name: "10000000000000000000000003"},
				{ID: 100004, Name: "10000000000000000000000004"},
				{ID: 100005, Name: "10000000000000000000000005"},
				{ID: 100006, Name: "10000000000000000000000006"},
			},
		},
		{
			ID: 2,
			Children: []*IdNode{{
				ID: 21,
			}, {ID: 22},
			},
		},
	}
	type ExportData struct {
		Uns    []*IdNode `json:"uns"`
		Labels []*IdNode `json:"labels"`
	}
	var exportData = ExportData{
		Uns:    nodes,
		Labels: []*IdNode{{Name: "关系 10000000000000000000000guanxi"}, {Name: "时序 100000000000000000000000seq"}},
	}
	bs, _ := json.Marshal(exportData)

	t.Log("json总大小:", len(bs))
	bigJson := base.NewReadCloserWrapper(bytes.NewBuffer(bs))

	treeNodeList := make([]*TreeNode, 0, len(nodes)*4)
	count := 0
	er := DecodeStreamedJson(bigJson, 16, 4,
		nodeGetChildren, node2tree, func(readSize int64, propName string, nodes []*TreeNode) {
			t.Logf("%s nodes[%d]: readSize:%d, %v\n", propName, len(nodes), readSize, nodes)
			treeNodeList = append(treeNodeList, nodes...)
			count += len(nodes)
		})
	t.Log("decode count = ", count)
	if er != nil {
		t.Error(er)
	}
	t.Logf("treeNodeList: %v\n", treeNodeList)
}

func TestStreamJsonArray(t *testing.T) {
	var nodes = []*IdNode{
		{
			ID: 1,
			Children: []*IdNode{{
				ID: 11,
			}, {ID: 12, Children: []*IdNode{{ID: 120001}, {ID: 120002, Children: []*IdNode{{ID: 12000201}, {ID: 12000202}}}}},

				{ID: 100003}, {ID: 100004}, {ID: 100005}, {ID: 100006},
			},
		},
		{
			ID: 2,
			Children: []*IdNode{{
				ID: 21,
			}, {ID: 22},
			},
		},
	}
	bs, _ := json.Marshal(nodes)
	bigJson := base.NewReadCloserWrapper(bytes.NewBuffer(bs))

	treeNodeList := make([]*TreeNode, 0, len(nodes)*4)
	count := 0

	er := DecodeStreamedJson(bigJson, 16, 4,
		nodeGetChildren, node2tree, func(readSize int64, propName string, nodes []*TreeNode) {
			t.Logf("%s nodes[%d]: readSize:%d, %v\n", propName, len(nodes), readSize, nodes)
			treeNodeList = append(treeNodeList, nodes...)
			count += len(nodes)
		})
	t.Log("decode count = ", count)
	if er != nil {
		t.Error(er)
	}
	t.Logf("treeNodeList: %v\n", treeNodeList)
	// 创建编码器
	output := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := NewStreamJsonEncoder(output, nodeGetChildren, nodeSetChildren, treeGetId, treeGetParentId, tree2node)
	// 分批处理数据
	parts := base.Partition(treeNodeList, 5)
	lastIndex := len(parts) - 1
	for i, batch := range parts {
		t.Logf("final? %v, batch[%d]: %v\n", i == lastIndex, i, batch)
		er := WriteBatch(encoder, batch, i == lastIndex)
		if er != nil {
			t.Fatal(er)
		}
	}
	t.Log("encodeJson:\n", output.String())
}
func TestPath(t *testing.T) {
	t.Log("dir:", filepath.Base(filepath.Dir("\\export\\export20251128140052_3938211876.json")))
	t.Log("Clean:", filepath.Base(constants.ExportRoot))
	//if filepath.Dir(paramFilePath) != filepath.Clean(constants.ExportRoot)
}
