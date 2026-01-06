package types

import (
	"backend/internal/common/constants"
	"backend/internal/common/utils/PathUtil"
	"sort"
)

func (t *TopicTreeResult) AddChild(child *TopicTreeResult) {
	if child == nil {
		return
	}
	child.ParentId = &t.Id
	t.Children = append(t.Children, child)
	SortUnsList(t.Children)
}
func (t *TopicTreeResult) GetCountChildren() int {
	if t.CountChildren == nil {
		count := 0
		if children := t.Children; len(children) > 0 {
			for _, child := range children {
				if child.PathType == constants.PathTypeFile {
					count++
				}
				count += child.GetCountChildren()
			}
		}
		t.CountChildren = &count
		t.HasChildren = count > 0
	}
	return *t.CountChildren
}
func SortUnsList(list []*TopicTreeResult) {
	sort.Sort(unsList(list))
}

type unsList []*TopicTreeResult

func (x unsList) Len() int { return len(x) }
func (x unsList) Less(i, j int) bool {
	return PathUtil.GetName(x[i].Path) < PathUtil.GetName(x[j].Path)
}
func (x unsList) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
