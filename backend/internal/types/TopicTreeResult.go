package types

import (
	"backend/internal/common/utils/PathUtil"
	"sort"
)

func (t *TopicTreeResult) AddChild(child *TopicTreeResult) {
	if child == nil {
		return
	}
	t.Children = append(t.Children, child)
	sort.Sort(unsList(t.Children))
}

type unsList []*TopicTreeResult

func (x unsList) Len() int { return len(x) }
func (x unsList) Less(i, j int) bool {
	return PathUtil.GetName(x[i].Path) < PathUtil.GetName(x[j].Path)
}
func (x unsList) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
