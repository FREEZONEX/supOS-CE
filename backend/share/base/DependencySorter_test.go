package base

import (
	"strconv"
	"testing"
)

type Node struct {
	id  int
	pid int
}

func (n *Node) GetId() int {
	return n.id
}
func (n *Node) GetParentId() int {
	return n.pid
}
func nn(id, pid int) *Node {
	return &Node{id: id, pid: pid}
}
func (n *Node) String() string {
	return strconv.Itoa(n.id)
}
func TestBuildReverseGraph(t *testing.T) {
	list := []*Node{nn(50, 50),
		//nn(50, 42), nn(42, 43), nn(43, 50),

		nn(31, 2), nn(30, 2), nn(40, 2), nn(2, 1)}
	sList, erList := SorByDependency(list, func(n *Node, n2 *Node) bool {
		xiao := n.id < n2.id
		t.Logf(" %d < %d ? %v\n", n.id, n2.id, xiao)
		return xiao
	}, func(n *Node) int {
		return n.id
	}, func(n *Node) int {
		return n.pid
	})
	t.Logf("erList: %s", erList)
	t.Logf("sList: %s", sList)
}
