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
	list := []*Node{nn(3, 2), nn(4, 2), nn(2, 1)}
	SorByDependency(list, func(n *Node) int {
		return n.id
	}, func(n *Node) int {
		return n.pid
	})
	t.Logf("levelMap: %s", list)

}
